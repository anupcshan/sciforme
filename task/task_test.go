package task_test

import (
	"fmt"
	"github.com/anupcshan/neoism"
	"github.com/anupcshan/sciforme/task"
	"github.com/jmcvetta/randutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"math"
	"testing"
)

const DESC_LEN = 20

type MockGraphDB struct {
	mock.Mock
}

func (m *MockGraphDB) CreateNode(props neoism.Props) (*neoism.Node, error) {
	args := m.Mock.Called(props)
	node, ok := args.Get(0).(*neoism.Node)
	if !ok {
		return nil, args.Error(1)
	}
	return node, args.Error(1)
}

func (m *MockGraphDB) NodesByLabel(label string) ([]*neoism.Node, error) {
	args := m.Mock.Called(label)
	nodes, ok := args.Get(0).([]*neoism.Node)
	if !ok {
		return nil, args.Error(1)
	}
	return nodes, args.Error(1)
}

func (m *MockGraphDB) Node(id int) (*neoism.Node, error) {
	args := m.Mock.Called(id)
	node, ok := args.Get(0).(*neoism.Node)
	if !ok {
		return nil, args.Error(1)
	}
	return node, args.Error(1)
}

func TestAddTaskNoDB(t *testing.T) {
	tm := task.TaskManager{}
	err, tsk := tm.AddTask("foo")

	assert.Error(t, err, "Error expected when DB instance was nil")
	assert.Nil(t, tsk, "Expected nil value for task")
}

func TestAddTaskDBNotConnected(t *testing.T) {
	mdb := new(MockGraphDB)
	tm := task.TaskManager{Database: mdb}

	str, _ := randutil.AlphaString(DESC_LEN)
	mdb.On("CreateNode", neoism.Props{"name": str}).Return(nil, fmt.Errorf(""))
	err, tsk := tm.AddTask(str)

	assert.Error(t, err, "Error expected when non-working DB instance was provided")
	assert.Nil(t, tsk, "Expected nil value for task")

	mdb.AssertExpectations(t)
}

func TestAddTaskSuccess(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	db, _ := neoism.Connect("http://localhost:7474/db/data")
	tm := task.TaskManager{Database: db}

	str, _ := randutil.AlphaString(DESC_LEN)
	err, tsk := tm.AddTask(str)

	assert.NoError(t, err, nil, "No error expected while adding a new node")
	assert.NotNil(t, tsk, "Expected non-nil value for task")

	node, err := db.Node(tsk.Id)
	assert.NoError(t, err, nil, "Could not find newly created node in DB")

	// Cleanup
	defer func() {
		node.Delete()
	}()

	name, err := node.Property("name")
	assert.NoError(t, err, nil, "'name' property not set")
	assert.Equal(t, name, str, "Testing task name")

	labels, _ := node.Labels()
	assert.True(t, len(labels) > 0, task.TASK_LABEL, "Testing task label")
	assert.Equal(t, labels[0], task.TASK_LABEL, "Testing task label")
}

func TestListTasksNoDB(t *testing.T) {
	tm := task.TaskManager{}
	err, tasks := tm.ListTasks()

	assert.Error(t, err, "Error expected when DB instance was nil")
	assert.Nil(t, tasks, "Expected nil value for list of tasks")
}

func TestListTasksDBNotConnected(t *testing.T) {
	mdb := new(MockGraphDB)
	tm := task.TaskManager{Database: mdb}

	mdb.On("NodesByLabel", task.TASK_LABEL).Return(nil, fmt.Errorf(""))
	err, tasks := tm.ListTasks()

	assert.Error(t, err, "Error expected when non-working DB instance was provided")
	assert.Nil(t, tasks, "Expected nil value for list of tasks")

	mdb.AssertExpectations(t)
}

func TestListTasksEmpty(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	db, _ := neoism.Connect("http://localhost:7474/db/data")
	tm := task.TaskManager{Database: db}

	err, tasks := tm.ListTasks()

	assert.NoError(t, err, nil, "No error expected while fetching empty list")
	assert.Empty(t, tasks, "Expected empty list of tasks")
}

func TestListTasksSingle(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	db, _ := neoism.Connect("http://localhost:7474/db/data")
	tm := task.TaskManager{Database: db}

	str, _ := randutil.AlphaString(DESC_LEN)
	err, tsk := tm.AddTask(str)
	node, err := db.Node(tsk.Id)

	// Cleanup
	defer func() {
		node.Delete()
	}()

	err, tasks := tm.ListTasks()

	assert.NoError(t, err, nil, "No error expected while listing tasks")
	assert.Equal(t, len(tasks), 1, "Expected list of size 1")
}

func TestAddDependencyNoDB(t *testing.T) {
	tm := task.TaskManager{}
	err, task := tm.AddDependency(1, 2)

	assert.Error(t, err, "Error expected when DB instance was nil")
	assert.Nil(t, task, "Expected nil value for list of tasks")
}

func TestAddDependencyDBNotConnected(t *testing.T) {
	mdb := new(MockGraphDB)
	tm := task.TaskManager{Database: mdb}

	id, _ := randutil.IntRange(0, math.MaxInt32)
	depId, _ := randutil.IntRange(0, math.MaxInt32)

	// Id error
	mdb.On("Node", id).Return(nil, fmt.Errorf(""))
	err, tsk := tm.AddDependency(id, depId)

	assert.Error(t, err, "Error expected when non-working DB instance was provided")
	assert.Nil(t, tsk, "Expected nil value for task")

	mdb.AssertExpectations(t)

	// DepId error
	mdb = new(MockGraphDB)
	tm = task.TaskManager{Database: mdb}

	mdb.On("Node", id).Return(nil, nil)
	mdb.On("Node", depId).Return(nil, fmt.Errorf(""))
	err, tsk = tm.AddDependency(id, depId)

	assert.Error(t, err, "Error expected when non-working DB instance was provided")
	assert.Nil(t, tsk, "Expected nil value for task")

	mdb.AssertExpectations(t)
}
