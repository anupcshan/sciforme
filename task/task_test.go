package task_test

import (
	"github.com/anupcshan/sciforme/task"
	"github.com/jmcvetta/neoism"
	"github.com/jmcvetta/randutil"
	"github.com/stretchr/testify/assert"
	"testing"
)

const DESC_LEN = 20

func TestAddTaskNoDB(t *testing.T) {
	tm := task.TaskManager{}
	err, tsk := tm.AddTask("foo")

	assert.Error(t, err, "Error expected when DB instance was nil")
	assert.Nil(t, tsk, "Expected nil value for task")
}

func TestAddTaskDBNotConnected(t *testing.T) {
	db, _ := neoism.Connect("http://localhost:12345/db/data")
	tm := task.TaskManager{Database: db}

	str, _ := randutil.AlphaString(DESC_LEN)
	err, tsk := tm.AddTask(str)

	assert.Error(t, err, "Error expected when non-working DB instance was provided")
	assert.Nil(t, tsk, "Expected nil value for task")
}

func TestAddTaskSuccess(t *testing.T) {
	db, _ := neoism.Connect("http://localhost:7474/db/data")
	tm := task.TaskManager{Database: db}

	str, _ := randutil.AlphaString(DESC_LEN)
	err, tsk := tm.AddTask(str)

	assert.NoError(t, err, nil, "No error expected while adding a new node")
	assert.NotNil(t, tsk, "Expected non-nil value for task")

	node, err := db.Node(tsk.Id)
	assert.NoError(t, err, nil, "Could not find newly created node in DB")

	// Cleanup
	defer node.Delete()

	name, err := node.Property("name")
	assert.NoError(t, err, nil, "'name' property not set")
	assert.Equal(t, name, str, "Testing task name")

	labels, _ := node.Labels()
	assert.Equal(t, labels[0], task.TASK_LABEL, "Testing task label")
}

func TestListTasksNoDB(t *testing.T) {
	tm := task.TaskManager{}
	err, tasks := tm.ListTasks()

	assert.Error(t, err, "Error expected when DB instance was nil")
	assert.Nil(t, tasks, "Expected nil value for list of tasks")
}

func TestListTasksDBNotConnected(t *testing.T) {
	db, _ := neoism.Connect("http://localhost:12345/db/data")
	tm := task.TaskManager{Database: db}

	err, tasks := tm.ListTasks()

	assert.Error(t, err, "Error expected when non-working DB instance was provided")
	assert.Nil(t, tasks, "Expected nil value for list of tasks")
}

func TestListTasksEmpty(t *testing.T) {
	db, _ := neoism.Connect("http://localhost:7474/db/data")
	tm := task.TaskManager{Database: db}

	err, tasks := tm.ListTasks()

	assert.NoError(t, err, nil, "No error expected while fetching empty list")
	assert.Empty(t, tasks, "Expected empty list of tasks")
}

func TestListTasksSingle(t *testing.T) {
	db, _ := neoism.Connect("http://localhost:7474/db/data")
	tm := task.TaskManager{Database: db}

	str, _ := randutil.AlphaString(DESC_LEN)
	err, tsk := tm.AddTask(str)
	node, err := db.Node(tsk.Id)
	// Cleanup
	defer node.Delete()

	err, tasks := tm.ListTasks()

	assert.NoError(t, err, nil, "No error expected while listing tasks")
	assert.Equal(t, len(tasks), 1, "Expected list of size 1")
}
