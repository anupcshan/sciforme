package task

import (
	"fmt"
	"github.com/golang/glog"
	"github.com/jmcvetta/neoism"
)

type Task struct {
}

type TaskManager struct {
	Database *neoism.Database
}

func (self TaskManager) AddTask(desc string) error {
	if self.Database == nil {
		return fmt.Errorf("No database connection defined")
	}

	td, err := self.Database.CreateNode(neoism.Props{"name": desc})
	td.AddLabel("Task")

	if err != nil {
		return fmt.Errorf("DB Error: %q", err)
	}

	if glog.V(2) {
		return fmt.Errorf("Added task: ", desc)
	}

	return nil
}
