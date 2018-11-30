package ecsupdatenotify

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecs"
)

// ListTasks ... Returns a list of tasks for a specified cluster.
func (m *Monitor) ListTasks() []*string {
	var taskArns []*string
	input := &ecs.ListTasksInput{
		Cluster: aws.String(m.Name),
	}

	m.Client.ecs.ListTasksPages(
		input,
		func(page *ecs.ListTasksOutput, _ bool) bool {
			for _, taskArn := range page.TaskArns {
				taskArns = append(taskArns, taskArn)
			}
			return true
		},
	)
	return taskArns
}

// DescribeTasks ... Returns a list of tasks for a specified cluster
func (m *Monitor) DescribeTasks(tasks []*string) []*ecs.Task {
	input := &ecs.DescribeTasksInput{
		Cluster: aws.String(m.Name),
		Tasks:   tasks,
	}

	result, err := m.Client.ecs.DescribeTasks(input)
	if err != nil {
		log.sugar.Warnf("failed to DescribeTasks: %s cluster\n", m.Name)
	}

	return result.Tasks
}

// DescribeTaskDefinition ... Describes a task definition
func (m *Monitor) DescribeTaskDefinition(taskDefinition *string) *ecs.DescribeTaskDefinitionOutput {
	input := &ecs.DescribeTaskDefinitionInput{
		TaskDefinition: taskDefinition,
	}

	result, err := m.Client.ecs.DescribeTaskDefinition(input)
	if err != nil {
		log.sugar.Warnf("failed to get image: %s cluster\n", m.Name)
	}

	return result
}

// NewClient ... Creates a new instance of the ECS client with a session.
func (m *Monitor) NewClient() *Client {
	sess, _ := session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
		Profile:           m.AWSProfile,
	})

	return &Client{
		ecs: ecs.New(sess, aws.NewConfig().WithRegion(m.AWSRegion)),
	}
}
