package main

import (
	"fmt"
	"github.com/anupcshan/neoism"
	"github.com/anupcshan/sciforme/task"
	"github.com/codegangsta/cli"
	"github.com/golang/glog"
	"os"
	"strconv"
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
		{
			Name:      "depends",
			ShortName: "dep",
			Usage:     "Add a dependency between 2 tasks",
			Action: func(c *cli.Context) {
				if !c.Args().Present() || len(c.Args()) < 2 {
					glog.Fatal("Depends: Please provide ids of 2 tasks to add dependency between")
				}

				id, err := strconv.Atoi(c.Args().Get(0))

				if err != nil {
					glog.Fatal("Could not decode argument %q into an integer, %q", c.Args().Get(0), err)
				}

				depId, err := strconv.Atoi(c.Args().Get(1))

				if err != nil {
					glog.Fatal("Could not decode argument %q into an integer, %q", c.Args().Get(1), err)
				}

				err, _ = tm.AddDependency(id, depId)

				if err != nil {
					glog.Fatal(err)
				}
			},
		},
	}

	app.Run(os.Args)
}
