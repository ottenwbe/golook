// +build integration

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
	//"encoding/json"
	"github.com/fsouza/go-dockerclient"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	//"github.com/ottenwbe/golook/broker/api"
	"github.com/ottenwbe/golook/broker/communication"
	"github.com/ottenwbe/golook/broker/models"
	"github.com/ottenwbe/golook/broker/routing"
	"github.com/ottenwbe/golook/broker/runtime"
	"github.com/ottenwbe/golook/broker/service"
	log "github.com/sirupsen/logrus"
	//"io/ioutil"
	//"net/http"
	"github.com/ottenwbe/golook/broker/utils"
	"time"
)

var _ = Describe("The golook server's", func() {

	/*Context("http api", func() {
		It("allows clients to query the application's /info endpoint from a server running in a docker container", func() {
			RunInDockerContainer(func(client *docker.Client, container *docker.Container) {

				appInfo := &runtime.AppInfo{}
				compareInfo := runtime.NewAppInfo()
				containerInfo := getContainerInfo(client, container)

				// make test request
				r, errGet := http.Get("http://" + containerInfo.NetworkSettings.IPAddress + ":8383" + api.INFO_EP)
				Expect(errGet).To(BeNil())

				// get the result
				b, errRead := ioutil.ReadAll(r.Body)
				Expect(errRead).To(BeNil())
				errMarshal := json.Unmarshal(b, appInfo)
				Expect(errMarshal).To(BeNil())

				// verify that the result is correct
				Expect(appInfo).ToNot(BeNil())
				Expect(appInfo.Version).To(Equal(compareInfo.Version))
				Expect(appInfo.System.IP).To(Equal(containerInfo.NetworkSettings.IPAddress))
			})
		})
	})*/
})

func getContainerInfo(client *docker.Client, container *docker.Container) *docker.Container {
	containerInfo, errC := client.InspectContainer(container.ID)
	if errC != nil {
		log.Fatal("Could not query for the container's info")
	}
	return containerInfo
}

var _ = Describe("The golook server's", func() {
	Context("rpc api", func() {
		It("registers a system via the rest api", func() {

			RunInDockerContainer(func(client *docker.Client, container *docker.Container) {

				var testRouter = routing.NewRouter("system")
				containerInfo, _ := client.InspectContainer(container.ID)
				testRpcClient := communication.NewLookupRPCClient("http://" + containerInfo.NetworkSettings.IPAddress + ":8382")
				testRouter.NewPeer(routing.NewKey(containerInfo.NetworkSettings.MacAddress), testRpcClient)

				testSystem := &models.SystemFiles{System: runtime.GolookSystem, Files: nil}
				res := testRouter.BroadCast(service.SYSTEM_REPORT,
					&service.SystemTransfer{Uuid: runtime.GolookSystem.UUID, System: testSystem, IsDeletion: false})

				actualResponse := service.PeerResponse{}
				err := utils.UnmarshalS(res, &actualResponse)

				Expect(err).To(BeZero())
				Expect(actualResponse.Success).To(BeTrue())
			})
		})
	})
})

func RunInDockerContainer(f func(client *docker.Client, container *docker.Container)) {
	client, err := docker.NewClientFromEnv()
	if err != nil {
		log.Fatalf("Cannot connect to Docker daemon: %s", err)
	}
	container, err := client.CreateContainer(createOptions("golook:latest"))
	if err != nil {
		log.Fatalf("Cannot create Docker container; make sure docker daemon is started: %s", err)
	}
	defer func() { //ensure container stops again
		if err := client.RemoveContainer(docker.RemoveContainerOptions{
			ID:    container.ID,
			Force: true,
		}); err != nil {
			log.Fatalf("cannot remove container: %s", err)
		}
	}()

	err = client.StartContainer(container.ID, &docker.HostConfig{})
	if err != nil {
		log.Fatalf("Cannot start Docker container: %s", err)
	}

	// some time for the container
	time.Sleep(time.Second)

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
