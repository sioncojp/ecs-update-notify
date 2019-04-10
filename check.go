package ecsupdatenotify

import (
	"strings"

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
	m.Client = m.NewClient()

	tasks := m.ListTasks()
	if len(tasks) == 0 {
		log.sugar.Warnf("task is not registered: %s cluster\n", m.Name)
	}

	describeTasks := m.DescribeTasks(tasks)

	// Check all task in cluster
	for _, v := range describeTasks {
		taskDefinition := strings.Split(*v.TaskDefinitionArn, "/")[1]
		task := strings.Split(taskDefinition, ":")[0]
		revision, _ := strconv.Atoi(strings.Split(taskDefinition, ":")[1])
		// initialize
		if len(m.Tasks) == 0 {
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
		m.CheckTasksUpdate(task, v.TaskDefinitionArn, revision, v.LastStatus)
	}

	// Notification when Revision is updated
	m.PostSlackMessage()
}

// CheckTasksUpdate ... Check all tasks and set value to struct
func (m *Monitor) CheckTasksUpdate(task string, taskDefinitionArn *string, revision int, l *string) {
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
