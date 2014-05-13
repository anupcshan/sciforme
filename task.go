package main

import (
	"github.com/codegangsta/cli"
	"github.com/jmcvetta/neoism"
  "github.com/golang/glog"
	"os"
)

func main() {
	db, err := neoism.Connect("http://localhost:7474/db/data")

	if err != nil {
		glog.Fatalf("DB Error: %q", err)
	}

	app := cli.NewApp()
	app.Name = "todo"
	app.Usage = "task list on the command line"
  app.Version = "0.1"
	app.Commands = []cli.Command{
		{
			Name:      "add",
			ShortName: "a",
			Usage:     "Add a task to the list",
			Action: func(c *cli.Context) {
        if !c.Args().Present() {
          glog.Error("Add: Please provide description of task to add")
        }

        tName := c.Args().First()

        td, err := db.CreateNode(neoism.Props{"name": tName})
        td.AddLabel("Task")

        if err != nil {
          glog.Fatalf("DB Error: %q", err)
        }

				if glog.V(2) {
          glog.Info("Added task: ", tName)
        }
			},
		},
		{
			Name:      "list",
			ShortName: "ls",
			Usage:     "List of tasks to be completed",
			Action: func(c *cli.Context) {
        glog.Fatal("TODO: list")
			},
		},
	}

	app.Run(os.Args)
}
