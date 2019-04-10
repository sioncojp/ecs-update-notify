package ecsupdatenotify

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"strings"

	"github.com/aws/aws-sdk-go/service/ecs"
)

type SlackMessage interface {
	NewAttachmentMessage(taskDefiniton *ecs.TaskDefinition, awsProfile, cluster, task string, revision int) *Attachment
}

// PostSlackMessage ... Verify the revision number and notify the message
func (m *Monitor) PostSlackMessage() {
	for _, t := range m.Tasks {
		td := m.DescribeTaskDefinition(t.TaskDefinitionArn).TaskDefinition

		if t.isCurrReivision {
			// notify failure
			if (t.FailureCount % CheckFailureInterval) == 0 {
				e := &ECSFailureTask{}
				failure := e.NewAttachmentMessage(td, m.AWSProfile, m.Name, t.Name, t.NextRevision)

				m.postSlackMessage(failure)
			}
			if t.NextRevision != t.CurrRevision {
				t.FailureCount++
				fmt.Println(t.FailureCount)
				break
			}
		} else {
			// notify success
			success := t.NewAttachmentMessage(td, m.AWSProfile, m.Name, t.Name, t.NextRevision)

			m.postSlackMessage(success)

			// reset
			t.CurrRevision = t.NextRevision
		}

		// reset
		t.isCurrReivision = false
		t.FailureCount = 1
	}
}

// postSlackMessage ... Message to slack
func (m *Monitor) postSlackMessage(a *Attachment) {
	s := Slack{
		Attachments: []*Attachment{a},
	}

	message, _ := json.Marshal(s)

	resp, _ := http.PostForm(
		m.IncomingWebhook,
		url.Values{"payload": {string(message)}},
	)

	if _, err := ioutil.ReadAll(resp.Body); err != nil {
		log.sugar.Warnf("failed to post message to slack: %s monitor: %s", m.Name, err)
	}
	defer resp.Body.Close()
}

// NewAttachmentMessage ... Initialize attachment data of slack for change task message
func (e *ECSTask) NewAttachmentMessage(taskDefinition *ecs.TaskDefinition, awsProfile, cluster, task string, revision int) *Attachment {
	var images []string

	cpu := *taskDefinition.Cpu
	memory := *taskDefinition.Memory

	for _, i := range taskDefinition.ContainerDefinitions {
		images = append(images, *i.Image)
	}

	return &Attachment{
		ColorORANGE,
		"ECS Update Notify",
		fmt.Sprintf("%s task updated: %d\n"+
			"CPU: %s, MEM: %sGB\n"+
			"Image: \n"+
			"```"+
			"%s"+
			"```", task, revision, cpu, memory, strings.Join(images, "\n")),
		fmt.Sprintf("%s cluster in %s", cluster, awsProfile),
	}
}

// NewAttachmentMessage ... Initialize attachment data of slack for failure messages
func (e *ECSFailureTask) NewAttachmentMessage(taskDefinition *ecs.TaskDefinition, awsProfile, cluster, task string, revision int) *Attachment {
	var images []string

	for _, i := range taskDefinition.ContainerDefinitions {
		images = append(images, *i.Image)
	}

	return &Attachment{
		ColorRED,
		"ECS Update Notify",
		fmt.Sprintf("If task failure???\n"+
			"Please check the ECS task\n"+
			"```"+
			"Task: %s:%d\n"+
			"Image: %s"+
			"```", task, revision, strings.Join(images, "\n")),
		fmt.Sprintf("%s cluster in %s", cluster, awsProfile),
	}
}
