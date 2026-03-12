package task_test

import (
	"testing"

	"github.com/example/scheduler/task"
	"github.com/stretchr/testify/assert"
)

func TestTaskKeyAndOrder(t *testing.T) {
	tsk := task.Task{Identifier: 42, Priority: task.HighPriority}
	assert.Equal(t, 42, tsk.Key())
	assert.Equal(t, int(task.HighPriority), tsk.Order())
}

func TestTaskIdentity(t *testing.T) {
	tsk := task.Task{Identifier: 7, Priority: task.UrgentPriority}
	assert.Equal(t, task.TaskID(7), tsk.Identifier)
	assert.Equal(t, task.UrgentPriority, tsk.Priority)
}
