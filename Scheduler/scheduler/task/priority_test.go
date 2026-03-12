package task_test

import (
	"testing"

	"github.com/example/scheduler/task"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewPriorityValid(t *testing.T) {
	p, err := task.NewPriority(50)
	require.NoError(t, err)
	assert.Equal(t, task.Priority(50), p)
}

func TestNewPriorityBoundaries(t *testing.T) {
	_, err := task.NewPriority(0)
	assert.NoError(t, err)

	_, err = task.NewPriority(100)
	assert.NoError(t, err)
}

func TestNewPriorityOutOfRange(t *testing.T) {
	_, err := task.NewPriority(-1)
	assert.Error(t, err)

	_, err = task.NewPriority(101)
	assert.Error(t, err)
}

func TestMustPriorityPanics(t *testing.T) {
	assert.Panics(t, func() { task.MustPriority(200) })
}

func TestMustPriorityValid(t *testing.T) {
	assert.NotPanics(t, func() { task.MustPriority(75) })
}

func TestPredefinedLevels(t *testing.T) {
	assert.Less(t, int(task.LowPriority), int(task.NormalPriority))
	assert.Less(t, int(task.NormalPriority), int(task.HighPriority))
	assert.Less(t, int(task.HighPriority), int(task.CriticalPriority))
	assert.Less(t, int(task.CriticalPriority), int(task.UrgentPriority))
}
