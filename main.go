package main

import (
	"fmt"
	"github.com/anupcshan/sciforme/task"
	"github.com/codegangsta/cli"
	"github.com/golang/glog"
	"github.com/jmcvetta/neoism"
	"os"
)

func main() {
	db, err := neoism.Connect("http://localhost:7474/db/data")

	if err != nil {
		glog.Fatalf("DB Error: %q", err)
	}

	tm := task.TaskManager{Database: db}

	app := cli.NewApp()
	app.Name = "sciforme"
	app.Usage = "Schedule it for me :: task list on the command line"
	app.Version = "0.1"
	app.Commands = []cli.Command{
		{
			Name:      "add",
			ShortName: "a",
			Usage:     "Add a task to the list",
			Action: func(c *cli.Context) {
				if !c.Args().Present() {
					glog.Fatal("Add: Please provide description of task to add")
				}

				tName := c.Args().First()
				tm.AddTask(tName)
			},
		},
		{
			Name:      "list",
			ShortName: "ls",
			Usage:     "List of tasks to be completed",
			Action: func(c *cli.Context) {
				err, list := tm.ListTasks()
				if err != nil {
					glog.Fatal("Error while fetching list of tasks", err)
				}

				if list != nil {
					for i := range list {
						fmt.Printf("%d:%s\n", list[i].Id, list[i].Name)
					}
				}
			},
		},
	}

	app.Run(os.Args)
}
