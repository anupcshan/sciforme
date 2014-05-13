package task_test

import (
	"github.com/anupcshan/sciforme/task"
	"reflect"
	"testing"
)

func expect(t *testing.T, a interface{}, b interface{}) {
	if a != b {
		t.Errorf("Expected %v (type %v) - Got %v (type %v)",
			b, reflect.TypeOf(b), a, reflect.TypeOf(a))
	}
}

func TestAddNoDB(t *testing.T) {
	tm := task.TaskManager{}
	err := tm.AddTask("foo")

	expect(t, err.Error(), "No database connection defined")
}
