package main

import (
	"fmt"
	"os"

	cli "github.com/jawher/mow.cli"
)

func main() {

	app := cli.App("cronr", "Executes scheduled jobs defined in a config file")

	// version
	app.Version("v version", "cronr 1 (build 1)")

	// options & args
	app.Spec = "-c [--verbose]"
	configFile := app.StringOpt("c config-file", "", "Path of the cronr config file")
	verbose := app.BoolOpt("verbose", false, "Display more informations on scheduled jobs")

	// action
	app.Action = func() {
		err := CronAction(*configFile, *verbose)
		if err != nil {
			fmt.Println(err)
			cli.Exit(1)
		}
	}

	app.Run(os.Args)
}
