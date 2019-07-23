package main

import "github.com/ggrcha/conductor-go-client"

func main() {
	c := conductor.NewConductorWorker("http://conductor-server:8080/api", 2, 10000)
	c.Start("backy_scrap", backyWorker, true)
}
