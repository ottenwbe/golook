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

package integration

import (
	"time"

	"github.com/fsouza/go-dockerclient"

	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"net/http"
)

// based on https://divan.github.io/posts/integration_testing/

func RunPeerInDocker(f func(client *docker.Client, container *docker.Container)) {
	var (
		client    *docker.Client
		container *docker.Container
		err       error
	)

	client, err = docker.NewClientFromEnv()
	failOnError(err, "Cannot connect to Docker daemon")

	container, err = client.CreateContainer(createOptions("golook:latest"))
	failOnError(err, "Cannot create Docker container; make sure docker daemon is started: %s")

	defer func() { //ensure container stops again
		if err := client.RemoveContainer(docker.RemoveContainerOptions{
			ID:    container.ID,
			Force: true,
		}); err != nil {
			log.Fatalf("cannot remove container: %s", err)
		}
	}()

	err = client.StartContainer(container.ID, &docker.HostConfig{})
	failOnError(err, "Cannot start Docker container")

	container, err = GetContainerInfo(client, container)
	failOnError(err, "Cannot inspect the container.")

	waitForGolook(container.NetworkSettings.IPAddress, 5*time.Second)

	f(client, container)
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

func GetContainerInfo(client *docker.Client, container *docker.Container) (information *docker.Container, err error) {
	// wait for container to wake up
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
	return fmt.Errorf("Cannot start container %s for %v", id, maxWait)
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
	return errors.New("Golook is not starting up in container.")
}

func failOnError(err error, message string) error {
	if err != nil {
		log.WithError(err).Fatal(message)
	}
	return err
}
