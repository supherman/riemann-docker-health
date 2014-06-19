package main

import (
	"fmt"
	"github.com/amir/raidman"
	"github.com/supherman/riemann-docker-health/docker"
	"github.com/supherman/riemann-docker-health/docker/cpu"
	"github.com/supherman/riemann-docker-health/docker/mem"
	"log"
	"os"
	"time"
)

func hostname() string {
	hostname, err := os.Hostname()
	if err != nil {
		log.Println(err)
	}
	return hostname
}

func alert(event *raidman.Event) {
	c, err := raidman.Dial("tcp", "localhost:5555")
	if err != nil {
		panic(err)
	}
	err = c.Send(event)

	if err != nil {
		panic(err)
	}
	c.Close()
}

func alertCPU(container string) {
	containerCpu, err := cpu.GetUsage(container)
	if err != nil {
		log.Println(err)
	}

	var cpuEvent = &raidman.Event{
		State:   "ok",
		Service: "cpu",
		Metric:  int(containerCpu),
		Ttl:     10,
		Host:    fmt.Sprintf("%s %s", hostname(), container),
	}
	alert(cpuEvent)
}

func alertMemory(container string) {
	containerMem, err := mem.GetUsage(container)
	if err != nil {
		log.Println(err)
	}

	var memEvent = &raidman.Event{
		State:   "ok",
		Service: "memory",
		Metric:  int(containerMem),
		Ttl:     10,
		Host:    fmt.Sprintf("%s %s", hostname(), container),
	}
	alert(memEvent)
}

func main() {
	go cpu.Monitor()
	tick := time.NewTicker(1 * time.Second)
	for {
		<-tick.C
		containers, err := docker.ListContainers()
		if err != nil {
			log.Println(err)
			continue
		}

		for _, container := range containers {
      alertCPU(container)
      alertMemory(container)
		}
	}
}
