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

func (self TaskManager) AddTask(desc string) (error, *Task) {
	if self.Database == nil {
		return fmt.Errorf("No database connection defined"), nil
	}

	td, err := self.Database.CreateNode(neoism.Props{"name": desc})
	td.AddLabel(TASK_LABEL)

	if err != nil {
		return fmt.Errorf("DB Error: %q", err), nil
	}

	if glog.V(2) {
		glog.Infoln("Added task: ", desc)
	}

	name, _ := td.Property("name")
	return nil, &Task{Name: name, Id: td.Id()}
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
