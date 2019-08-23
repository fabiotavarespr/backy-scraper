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

func isValidTask(configMap map[string]string) bool {
	_, poolOk := configMap["pool"]
	_, nameOk := configMap["image"]
	return poolOk && nameOk
}

func getPropertiesFromWorkerConfig(t *task.Task) (map[string]string, error) {
	w, _ := t.InputData["workerConfig"].(string)
	configs := strings.Split(w, ";")
	configMap := make(map[string]string)
	for _, config := range configs {
		keyValue := strings.SplitN(config, "=", 2)
		configMap[keyValue[0]] = keyValue[1]
	}
	return configMap, nil
}

func backyWorkerBackup(t *task.Task) (taskResult *task.TaskResult, err error) {

	// pool and name
	var p, n string

	// output data
	od := make(map[string]interface{})

	configMap, err := getPropertiesFromWorkerConfig(t)
	if err != nil {
		err = errors.New("could't parse workerConfig paramenters, must be key=value;key=value and is " + t.InputData["workerConfig"].(string))
		taskResult = task.NewTaskResult(t)
		taskResult.Status = "FAILED"
		taskResult.OutputData = od
		return taskResult, err
	}

	// retrieves input data
	if !isValidTask(configMap) {
		err = errors.New("invalid input parameters")
		taskResult = task.NewTaskResult(t)
		taskResult.Status = "FAILED"
		taskResult.OutputData = od
		return taskResult, err
	}

	// p, _ = t.InputData["pool"].(string)
	p, _ = configMap["pool"]
	// n, _ = t.InputData["name"].(string)
	n, _ = configMap["image"]

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

func backyWorkerRemove(t *task.Task) (taskResult *task.TaskResult, err error) {

	// dataId
	var id string

	// output data
	od := make(map[string]interface{})

	configMap, err := getPropertiesFromWorkerConfig(t)
	if err != nil {
		err = errors.New("could't parse workerConfig paramenters, must be key=value;key=value and is " + t.InputData["workerConfig"].(string))
		taskResult = task.NewTaskResult(t)
		taskResult.Status = "FAILED"
		taskResult.OutputData = od
		return taskResult, err
	}

	// retrieves input data
	if !isValidTask(configMap) {
		err = errors.New("invalid input parameters")
		taskResult = task.NewTaskResult(t)
		taskResult.Status = "FAILED"
		taskResult.OutputData = od
		return taskResult, err
	}

	id, _ = t.InputData["dataId"].(string)

	logrus.Infof("task-id: %s / dataId: %s", t.TaskId, id)

	id, lErr := doRemoveCleanup(id)

	if lErr != nil {
		err = lErr
		taskResult = task.NewTaskResult(t)
		taskResult.Status = "FAILED"
		taskResult.OutputData = od
		return taskResult, err
	}

	taskResult = task.NewTaskResult(t)
	taskResult.Status = "COMPLETED"
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

func doRemoveCleanup(id string) (string, error) {

	// Remove Backup command..
	backyCmd := cmd.NewCmd("backy2", "rm", id)
	statusChan := backyCmd.Start()

	finalStatus := <-statusChan

	log.Println("diff-bkp out: ", finalStatus.Stdout)
	log.Println("diff-bkp err: ", finalStatus.Stderr)

	if finalStatus.Exit != 0 {
		return "", (errors.New("backup remove failed"))
	}

	// Cleanup backy2 at selected object storage
	backyCmd2 := cmd.NewCmd("backy2", "cleanup")
	statusChan2 := backyCmd2.Start()

	finalStatus2 := <-statusChan2

	log.Println("diff-bkp out: ", finalStatus2.Stdout)
	log.Println("diff-bkp err: ", finalStatus2.Stderr)

	if finalStatus2.Exit != 0 {
		return "", (errors.New("backup cleanup failed"))
	}

	return "Backup removed and cleaned", nil
}
