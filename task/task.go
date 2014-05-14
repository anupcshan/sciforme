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

const ERR_TASK = -1

func (self TaskManager) AddTask(desc string) (error, int) {
	if self.Database == nil {
		return fmt.Errorf("No database connection defined"), ERR_TASK
	}

	td, err := self.Database.CreateNode(neoism.Props{"name": desc})
	td.AddLabel("Task")

	if err != nil {
		return fmt.Errorf("DB Error: %q", err), ERR_TASK
	}

	if glog.V(2) {
		return fmt.Errorf("Added task: ", desc), ERR_TASK
	}

	return nil, td.Id()
}
