package main

import (
	"fmt"
	"github.com/samalba/dockerclient"
	"github.com/davecheney/mdns"
	"log"
	"regexp"
	"strings"
)

// Callback used to listen to Docker's events
func eventCallback(event *dockerclient.Event, ec chan error, args ...interface{}) {
	if "start" == event.Status || "unpause" == event.Status || "die" == event.Status || "pause" == event.Status {
		docker, _ := dockerclient.NewDockerClient("unix:///var/run/docker.sock", nil)
		info, _ := docker.InspectContainer(event.ID)
		name := sanitize_domain_name(info.Name)
		ip := info.NetworkSettings.IPAddress

		run, _ := regexp.MatchString("_run_[[:digit:]]+$", info.Name)
		if run {
			return
		}

		if "start" == event.Status || "unpause" == event.Status {
			log.Printf("%s 300 IN A %s", name, ip)
			mdns.Publish(fmt.Sprintf("%s 300 IN A %s", name, ip))
		}
	}
}

func sanitize_domain_name(name string) string {
	return strings.Replace(strings.TrimLeft(name, "/"), "_", "-", -1) + ".local."
}

func main() {
	// Init the client
	docker, _ := dockerclient.NewDockerClient("unix:///var/run/docker.sock", nil)

	// Listen to events
	docker.StartMonitorEvents(eventCallback, nil)

	// Hold the execution to look at the events coming
	select{}
}
