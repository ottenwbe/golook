//Copyright 2016-2017 Beate Ottenwälder
//
//Licensed under the Apache License, Version 2.0 (the "License");
//you may not use this file except in compliance with the License.
//You may obtain a copy of the License at
//
//http://www.apache.org/licenses/LICENSE-2.0
//
//Unless required by applicable law or agreed to in writing, software
//distributed under the License is distributed on an "AS IS" BASIS,
//WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//See the License for the specific language governing permissions and
//limitations under the License.

/*
The docker integration allows us to start N docker images with a running golook broker, e.g., for testing.
It is based on https://divan.github.io/posts/integration_testing/
*/

package integration

import (
	"errors"
	"fmt"
	"github.com/fsouza/go-dockerclient"
	log "github.com/sirupsen/logrus"
	"net/http"
	"time"
)

/*
RunPeerInDocker starts one docker container with a running golook broker and then executes the function 'f'
*/
func RunPeerInDocker(f func(client *docker.Client, container *docker.Container)) {
	var (
		d = &DockerizedGolook{}
	)

	d.Init()
	//ensure Container stops again
	defer d.Stop()

	d.Start()

	f(d.Client, d.Container)
}

/*
RunPeersInDocker starts 'numPeers' docker containers with a running golook broker and then executes the function 'f'
*/
func RunPeersInDocker(numPeers int, f func(client []*DockerizedGolook)) {
	var (
		dockerizedGolooks = make([]*DockerizedGolook, numPeers)
	)

	for i := 0; i < numPeers; i++ {
		dockerizedGolooks[i] = &DockerizedGolook{}
		dockerizedGolooks[i].Init()
		//ensure Container stops again
		dockerizedGolooks[i].Start()
	}
	defer func() {
		for i := range dockerizedGolooks {
			dockerizedGolooks[i].Stop()
		}
	}()

	f(dockerizedGolooks)
}

/*
DockerizedGolook is a representation of a golook app running in a docker container
*/
type DockerizedGolook struct {
	Client    *docker.Client
	Container *docker.Container
}

/*
Init the docker container running the golook app
*/
func (d *DockerizedGolook) Init() {
	var err error

	d.Client, err = docker.NewClientFromEnv()
	failOnError(err, "Cannot connect to Docker daemon")

	d.Container, err = d.Client.CreateContainer(createOptions("golook:latest"))
	failOnError(err, "Cannot create Docker Container; make sure docker daemon is started: %s")

}

/*
Start the docker container running the golook app
*/
func (d *DockerizedGolook) Start() {
	var err error
	err = d.Client.StartContainer(d.Container.ID, &docker.HostConfig{})
	failOnError(err, "Cannot start Docker Container")

	d.Container, err = GetContainerInfo(d.Client, d.Container)
	failOnError(err, "Cannot inspect the Container.")

	waitForGolook(d.Container.NetworkSettings.IPAddress, 5*time.Second)

}

/*
Stop the docker container running the golook app
*/
func (d *DockerizedGolook) Stop() {
	if d.Client != nil {
		if err := d.Client.RemoveContainer(docker.RemoveContainerOptions{
			ID:    d.Container.ID,
			Force: true,
		}); err != nil {
			log.Fatalf("Cannot remove Container: %s", err)
		}
	} else {
		log.Fatal("Nil Container")
	}
}

func createOptions(containerName string) docker.CreateContainerOptions {
	ports := make(map[docker.Port]struct{})
	ports["8383"] = struct{}{}
	ports["8382"] = struct{}{}
	opts := docker.CreateContainerOptions{
		Config: &docker.Config{
			Image:        containerName,
			ExposedPorts: ports,
		},
	}

	return opts
}

/*
GetContainerInfo returns container information
*/
func GetContainerInfo(client *docker.Client, container *docker.Container) (information *docker.Container, err error) {
	// wait for Container to wake up
	err = waitForDocker(client, container.ID, 5*time.Second)
	if err != nil {
		return nil, err
	}

	information, err = client.InspectContainer(container.ID)

	return information, err
}

func waitForDocker(client *docker.Client, id string, maxWait time.Duration) error {
	done := time.Now().Add(maxWait)
	for time.Now().Before(done) {
		c, err := client.InspectContainer(id)
		if err != nil {
			break
		}
		if c.State.Running {
			return nil
		}
		time.Sleep(100 * time.Millisecond)
	}
	return fmt.Errorf("Cannot start Container %s for %v", id, maxWait)
}

func waitForGolook(ip string, maxWait time.Duration) error {
	done := time.Now().Add(maxWait)
	for time.Now().Before(done) {
		r, _ := http.Get("http://" + ip + ":8383/")

		if r != nil {
			return nil
		}
		time.Sleep(100 * time.Millisecond)
	}
	return errors.New("Golook is not starting up in Container.")
}

func failOnError(err error, message string) error {
	if err != nil {
		log.WithError(err).Fatal(message)
	}
	return err
}
