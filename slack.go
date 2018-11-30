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

// PostSlackMessage ... Verify the revision number and notify the message
func (m *Monitor) PostSlackMessage() {
	for _, t := range m.Tasks {
		if !t.isCurrReivision {
			td := m.DescribeTaskDefinition(t.TaskDefinitionArn).TaskDefinition
			a := t.NewAttachmentMessage(td, m.Name, t.Name, t.NextRevision)

			// post slack
			m.postSlackMessage(a)

			// reset
			t.CurrRevision = t.NextRevision
		}
		// reset
		t.isCurrReivision = false
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

// NewAttachmentMessage ... Initialize attachment data of slack
func (e *ECSTask) NewAttachmentMessage(taskDefiniton *ecs.TaskDefinition, cluster, task string, revision int) *Attachment {
	var images []string

	cpu := *taskDefiniton.Cpu
	memory := *taskDefiniton.Memory

	for _, i := range taskDefiniton.ContainerDefinitions {
		images = append(images, *i.Image)
	}

	return &Attachment{
		"#F6D64F",
		"ECS Update Notify",
		fmt.Sprintf("%s task updated: %d\n"+
			"CPU: %s, MEM: %sGB\n"+
			"Image: \n"+
			"```"+
			"%s"+
			"```", task, revision, cpu, memory, strings.Join(images, "\n")),
		fmt.Sprintf("%s cluster", cluster),
	}
}
