package task

import (
	"fmt"
	"github.com/anupcshan/neoism"
	"github.com/golang/glog"
)

type Task struct {
	Id           int
	Name         string
	dependencies []*Task
}

type TaskManager struct {
	Database neoism.GraphDB
}

const ERR_TASK = -1
const TASK_LABEL = "Task"
const DEPENDS_ON = "DependsOn"

func (self TaskManager) AddTask(desc string) (error, *Task) {
	if self.Database == nil {
		return fmt.Errorf("No database connection defined"), nil
	}

	td, err := self.Database.CreateNode(neoism.Props{"name": desc})

	if err != nil {
		return fmt.Errorf("DB Error: %q", err), nil
	}

	td.AddLabel(TASK_LABEL)

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

func (self TaskManager) AddDependency(id int, depId int) (error, *Task) {
	if self.Database == nil {
		return fmt.Errorf("No database connection defined"), nil
	}

	node, err := self.Database.Node(id)

	if err != nil {
		return fmt.Errorf("Could not find node with id %d, %q", id, err), nil
	}

	_, err = self.Database.Node(depId)

	if err != nil {
		return fmt.Errorf("Could not find node with id %d, %q", depId, err), nil
	}

	_, err = node.Relate(DEPENDS_ON, depId, neoism.Props{})

	if err != nil {
		return fmt.Errorf("Could not create relationship: %q", err), nil
	}

	node, err = self.Database.Node(id)

	if err != nil {
		return fmt.Errorf("Could not find node with id %d, %q", id, err), nil
	}

	return nil, nil
}
