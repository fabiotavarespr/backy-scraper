package main

import (
	"os"

	"github.com/ggrcha/conductor-go-client"
)

func main() {
	c := conductor.NewConductorWorker("http://"+os.Getenv("CONDUCTOR_HOST")+":"+os.Getenv("CONDUCTOR_PORT")+"/api", 2, 10000)
	typeWork := os.Getenv("CONDUCTOR_WORKNAME")
	switch typeWork {
	case "backup":
		c.Start(os.Getenv("CONDUCTOR_WORKNAME"), backyWorkerBackup, true)
	case "remove":
		c.Start(os.Getenv("CONDUCTOR_WORKNAME"), backyWorkerRemove, true)
	}
}
