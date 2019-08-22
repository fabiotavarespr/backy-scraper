package main

import (
	"os"

	"github.com/ggrcha/conductor-go-client"
)

func main() {
	c := conductor.NewConductorWorker("http://"+os.Getenv("CONDUCTOR_HOST")+":"+os.Getenv("CONDUCTOR_PORT")+"/api", 2, 10000)
	c.Start("backy_scrap", backyWorker, true)
}
