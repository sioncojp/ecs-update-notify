package ecsupdatenotify

import (
	"strings"

	"github.com/aws/aws-sdk-go/service/ecs"

	"github.com/aws/aws-sdk-go/aws"

	"strconv"

	"sync"
)

// CheckUpdate ... Check if it is replaced with a new task.
func (c *Config) CheckUpdate() {
	wg := &sync.WaitGroup{}

	// 10 parallel
	semaphore := make(chan struct{}, 10)

	for _, m := range c.Monitors {
		wg.Add(1)
		semaphore <- struct{}{}

		go func(m *Monitor) {
			defer func() {
				<-semaphore
				defer wg.Done()
			}()
			m.CheckClusterUpdate()
		}(m)
	}
	wg.Wait()
}

// CheckClusterUpdate ... Check update of task in cluster
func (m *Monitor) CheckClusterUpdate() {
	// initialize for nil pointer
	if len(m.Tasks) == 0 {
		m.Tasks = []*ECSTask{}
	}

	m.Client = m.NewClient()

	tasks := m.ListTasks()
	if len(tasks) == 0 {
		log.sugar.Warnf("task is not registered: %s cluster\n", m.Name)
	}

	d := m.DescribeTasks(tasks)
	describeTasks := make([]*ecs.Task, 0)

	for _, v := range d {
		// "arn:aws:ecs:ap-northeast-1:123456789:task-definition/xxxxx:1"
		// Extracting data xxxxxx
		if !ContainsString(
			m.IgnoreTasks,
			strings.Split(strings.Split(aws.StringValue(v.TaskDefinitionArn), "/")[1], ":")[0],
		) {
			describeTasks = append(describeTasks, v)
		}
	}

	// Check all task in cluster
	for _, v := range describeTasks {
		taskDefinition := strings.Split(*v.TaskDefinitionArn, "/")[1]
		task := strings.Split(taskDefinition, ":")[0]
		revision, _ := strconv.Atoi(strings.Split(taskDefinition, ":")[1])

		// first initialize
		if !IsContainsTaskName(task, m.Tasks) {
			m.Tasks = append(m.Tasks, &ECSTask{
				Name:              task,
				CurrRevision:      revision,
				NextRevision:      revision,
				isCurrReivision:   false,
				TaskDefinitionArn: v.TaskDefinitionArn,
				FailureCount:      1,
			})
		}

		// Check all tasks and set value to struct
		m.CheckTasksUpdate(task, v.TaskDefinitionArn, revision)
	}

	// Notification when Revision is updated
	m.PostSlackMessage()
}

// CheckTasksUpdate ... Check all tasks and set value to struct
func (m *Monitor) CheckTasksUpdate(task string, taskDefinitionArn *string, revision int) {
	for _, t := range m.Tasks {
		if t.Name == task {
			if t.CurrRevision == revision {
				t.isCurrReivision = true
			} else {
				t.NextRevision = revision
				t.TaskDefinitionArn = taskDefinitionArn
			}
		}
	}
}

// IsContainsTaskName ... ECSTask contains the specified task name
func IsContainsTaskName(task string, tasks []*ECSTask) bool {
	for _, t := range tasks {
		if task == t.Name {
			return true
		}
	}
	return false
}

// ContainsString...check to include string in the slice
func ContainsString(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
