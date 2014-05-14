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
	err, id := tm.AddTask("foo")

	assert.Error(t, err, "Error expected when DB instance was nil")
	assert.Equal(t, id, task.ERR_TASK, "Expected error value for new task id")
}

func TestAddTaskDBNotConnected(t *testing.T) {
	db, _ := neoism.Connect("http://localhost:12345/db/data")
	tm := task.TaskManager{Database: db}

	str, _ := randutil.AlphaString(DESC_LEN)
	err, id := tm.AddTask(str)

	assert.Error(t, err, "Error expected when non-working DB instance was provided")
	assert.Equal(t, id, task.ERR_TASK, "Expected error value for new task id")
}

func TestAddTaskSuccess(t *testing.T) {
	db, _ := neoism.Connect("http://localhost:7474/db/data")
	tm := task.TaskManager{Database: db}

	str, _ := randutil.AlphaString(DESC_LEN)
	err, id := tm.AddTask(str)

	assert.NoError(t, err, nil, "No error expected while adding a new node")
	assert.NotEqual(t, id, -1, "Expected non-error value for new task id")

	node, err := db.Node(id)
	assert.NoError(t, err, nil, "Could not find newly created node in DB")

	// Cleanup
	defer node.Delete()

	name, err := node.Property("name")
	assert.NoError(t, err, nil, "'name' property not set")
	assert.Equal(t, name, str, "Testing task name")

	labels, _ := node.Labels()
	assert.Equal(t, labels[0], task.TASK_LABEL, "Testing task label")
}
