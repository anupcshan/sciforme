package task_test

import (
	"github.com/anupcshan/sciforme/task"
	"github.com/jmcvetta/neoism"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAddTaskNoDB(t *testing.T) {
	tm := task.TaskManager{}
	err, id := tm.AddTask("foo")

	assert.Error(t, err, "Error expected when DB instance was nil")
	assert.Equal(t, id, task.ERR_TASK, "Expected error value for new task id")
}

func TestAddTaskDBNotConnected(t *testing.T) {
	db, _ := neoism.Connect("http://localhost:12345/db/data")
	tm := task.TaskManager{Database: db}
	err, id := tm.AddTask("foo")

	assert.Error(t, err, "Error expected when non-working DB instance was provided")
	assert.Equal(t, id, task.ERR_TASK, "Expected error value for new task id")
}

func TestAddTaskSuccess(t *testing.T) {
	// Test data not currently being cleaned up.
	db, _ := neoism.Connect("http://localhost:7474/db/data")
	tm := task.TaskManager{Database: db}
	err, id := tm.AddTask("foo")

	// Check if data exists in DB.
	assert.NoError(t, err, nil, "No error expected while adding a new node")
	assert.NotEqual(t, id, -1, "Expected non-error value for new task id")
}
