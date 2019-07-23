package main

import (
	"errors"
	"log"
	"regexp"
	"strings"

	"github.com/ggrcha/conductor-go-client/task"
	"github.com/go-cmd/cmd"
	"github.com/sirupsen/logrus"
)

func isValidTask(t *task.Task) bool {
	_, poolOk := t.InputData["pool"]
	_, nameOk := t.InputData["name"]
	return poolOk && nameOk
}

func backyWorker(t *task.Task) (taskResult *task.TaskResult, err error) {

	// pool and name
	var p, n string

	// output data
	od := make(map[string]interface{})

	// retrieves input data
	if !isValidTask(t) {
		err = errors.New("invalid input parameters")
		taskResult = task.NewTaskResult(t)
		taskResult.Status = "FAILED"
		taskResult.OutputData = od
		return taskResult, err
	}

	p, _ = t.InputData["pool"].(string)
	n, _ = t.InputData["name"].(string)

	logrus.Infof("task-id: %s / pool: %s", t.TaskId, p)
	logrus.Infof("task-id: %s / image: %s", t.TaskId, n)

	id, lErr := doBackup(p, n)

	if lErr != nil {
		err = lErr
		taskResult = task.NewTaskResult(t)
		taskResult.Status = "FAILED"
		taskResult.OutputData = od
		return taskResult, err
	}

	od["id"] = id

	taskResult = task.NewTaskResult(t)
	taskResult.Status = "COMPLETED"
	taskResult.OutputData = od
	err = nil

	return taskResult, err
}

func doBackup(p, i string) (string, error) {

	backyCmd := cmd.NewCmd("bash", "diff-bkp.sh", p, i)
	statusChan := backyCmd.Start()

	finalStatus := <-statusChan

	log.Println("diff-bkp out: ", finalStatus.Stdout)
	log.Println("diff-bkp err: ", finalStatus.Stderr)

	if finalStatus.Exit != 0 {
		return "", (errors.New("backup failed"))
	}

	rex, err := regexp.Compile("New version\\: ([\\-a-z0-9]+) \\(Tags")
	if err != nil {
		return "", errors.New(strings.Join(finalStatus.Stdout, "\r\n"))
	}
	id := rex.FindStringSubmatch(strings.Join(finalStatus.Stdout, "\r\n"))

	return id[1], nil
}
