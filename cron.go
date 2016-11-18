package main

import (
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"os/signal"

	"github.com/robfig/cron"
	yaml "gopkg.in/yaml.v2"
)

type command struct {
	name string
	args []string
}

// CronAction parses and schedules jobs, waiting for SIGINT signal to stop
func CronAction(configFilepath string, verbose bool) error {

	// get config file content
	fileContent, err := ioutil.ReadFile(configFilepath)
	if err != nil {
		return err
	}

	// parse config file
	config := ConfigFile{}
	err = yaml.Unmarshal(fileContent, &config)
	if err != nil {
		return err
	}

	// prepare cron
	c := cron.New()

	// creating jobs
	for i := range config.Jobs {
		job := config.Jobs[i]

		c.AddFunc(job.Cron, func() {

			if verbose {
				log.Println(job)
			}

			cmdArgs := newCommand(job.Command)
			cmd := exec.Command(cmdArgs.name, cmdArgs.args...)

			if job.WorkingDir != "" {
				cmd.Dir = job.WorkingDir
			}

			// if job.User != 0 && job.Group != 0 {
			// 	fmt.Println(job.User, ":", job.Group)
			// 	cmd.SysProcAttr = &syscall.SysProcAttr{}
			// 	cmd.SysProcAttr.Credential = &syscall.Credential{Uid: job.User, Gid: job.Group}
			// }

			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr

			err := cmd.Run()
			if err != nil {
				log.Printf("[ERROR] cmd: %v: %v\n", job.Command, err)
			}

		})
	}

	log.Printf("%v job(s) scheduled\n", len(c.Entries()))

	// start cron
	c.Start()

	// prepare chan for catching signals
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, os.Kill)

	// waiting for signal to stop
	<-done

	c.Stop()

	return nil
}

func newCommand(rawCommand []string) command {

	var (
		name string
		args []string
	)

	if len(rawCommand) > 1 {
		name = rawCommand[0]
		args = rawCommand[1:]
	} else {
		name = rawCommand[0]
		args = []string{}
	}

	return command{name: name, args: args}
}
