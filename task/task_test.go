package task_test

import (
	"github.com/anupcshan/sciforme/task"
	"github.com/jmcvetta/neoism"
	"reflect"
	"testing"
)

func expect(t *testing.T, a interface{}, b interface{}) {
	if a != b {
		t.Errorf("Expected %v (type %v) - Got %v (type %v)",
			b, reflect.TypeOf(b), a, reflect.TypeOf(a))
	}
}

func TestAddTaskNoDB(t *testing.T) {
	tm := task.TaskManager{}
	err := tm.AddTask("foo")

	expect(t, err.Error(), "No database connection defined")
}

func TestAddTaskDBNotConnected(t *testing.T) {
	db, _ := neoism.Connect("http://localhost:12345/db/data")
	tm := task.TaskManager{Database: db}
	err := tm.AddTask("foo")

	expect(t, err.Error(), "No database connection defined")
}

func TestAddTaskSuccess(t *testing.T) {
	// Test data not currently being cleaned up.
	db, _ := neoism.Connect("http://localhost:7474/db/data")
	tm := task.TaskManager{Database: db}
	err := tm.AddTask("foo")

	// Check if data exists in DB.

	expect(t, err, nil)
}
