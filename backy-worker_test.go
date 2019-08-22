package main

import (
	"reflect"
	"testing"

	"github.com/ggrcha/conductor-go-client/task"
	"github.com/sirupsen/logrus"
	"gotest.tools/assert"
)

func TestConfigWorkerParse(t *testing.T) {
	task := task.NewTask()
	task.InputData["workerConfig"] = "pool=volumes;image=teste-1=2"
	res, err := getPropertiesFromWorkerConfig(task)
	if err != nil {
		panic("Error at test getPropertiesFromWorkerConfig")
	}
	assert.Assert(t, reflect.DeepEqual(res, map[string]string{"image": "teste-1=2", "pool": "volumes"}))
	logrus.Infoln("Result")
	logrus.Infoln(res)
}

func TestIsValidConfig(t *testing.T) {
	configBase := map[string]string{"image": "teste-1=2", "pool": "volumes"}
	res := isValidTask(configBase)
	logrus.Info(res)
}
