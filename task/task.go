package task

import (
	"fmt"
	"github.com/golang/glog"
	"github.com/jmcvetta/neoism"
)

type Task struct {
	Id           int
	Name         string
	dependencies []*Task
}

type TaskManager struct {
	Database *neoism.Database
}

const ERR_TASK = -1
const TASK_LABEL = "Task"

func (self TaskManager) AddTask(desc string) (error, int) {
	if self.Database == nil {
		return fmt.Errorf("No database connection defined"), ERR_TASK
	}

	td, err := self.Database.CreateNode(neoism.Props{"name": desc})
	td.AddLabel(TASK_LABEL)

	if err != nil {
		return fmt.Errorf("DB Error: %q", err), ERR_TASK
	}

	if glog.V(2) {
		glog.Infoln("Added task: ", desc)
	}

	return nil, td.Id()
}

func (self TaskManager) ListTasks() (error, []*Task) {
	if self.Database == nil {
		return fmt.Errorf("No database connection defined"), nil
	}

	nodes, err := self.Database.NodesByLabel(TASK_LABEL)

	if err != nil {
		return fmt.Errorf("DB Error: %q", err), nil
	}

	var tasks []*Task

	for nId := range nodes {
		name, _ := nodes[nId].Property("name")
		tasks = append(tasks, &Task{Name: name, Id: nodes[nId].Id()})
	}

	return nil, tasks
}
