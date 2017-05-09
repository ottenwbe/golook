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

package service

import (
	"fmt"
	"github.com/fsouza/go-dockerclient"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/ottenwbe/golook/broker/models"
	"github.com/ottenwbe/golook/broker/repository"
	"github.com/ottenwbe/golook/broker/routing"
	golook "github.com/ottenwbe/golook/broker/runtime/core"
	"github.com/ottenwbe/golook/client"
	int "github.com/ottenwbe/golook/test/integration"
	"github.com/ottenwbe/golook/utils"
	log "github.com/sirupsen/logrus"
	"path/filepath"
	"time"
)

var _ = Describe("The service layer's", func() {

	var (
		systemService *SystemService
		fileServices  FileServices
		routerSystem  *router
		routerFiles   *router
	)

	BeforeEach(func() {
		fileServices = OpenFileServices(BroadcastFiles)
		routerFiles = fileServices.(*scenarioBroadcastFiles).r

		repositories.GoLookRepository = repositories.NewRepository()

		systemService = newSystemService()
		routerSystem = systemService.router
	})

	AfterEach(func() {
		repositories.GoLookRepository = repositories.NewRepository()
		systemService.close()
		CloseFileServices(fileServices)
	})

	Context("broadcast query service", func() {
		It("informs ALL peers about files in folders", func() {
			int.RunPeersInDocker(2, func(dockerizedGolook []*int.DockerizedGolook) {

				d1 := dockerizedGolook[0]
				d2 := dockerizedGolook[1]

				routerFiles.NewPeer(routing.NewKey(d1.Container.NetworkSettings.IPAddress), d1.Container.NetworkSettings.IPAddress)
				routerSystem.NewPeer(routing.NewKey(d1.Container.NetworkSettings.IPAddress), d1.Container.NetworkSettings.IPAddress)
				routerFiles.NewPeer(routing.NewKey(d2.Container.NetworkSettings.IPAddress), d2.Container.NetworkSettings.IPAddress)
				routerSystem.NewPeer(routing.NewKey(d2.Container.NetworkSettings.IPAddress), d2.Container.NetworkSettings.IPAddress)

				systemService.Run()
				_, err := fileServices.Report(&models.FileReport{Path: "integration_test.go"})
				if err != nil {
					log.Fatal("Cannot make file report in test.")
				}

				time.Sleep(time.Second)

				// query both of the peers for the distributed file
				client.Host = fmt.Sprintf("http://%s:8383", d1.Container.NetworkSettings.IPAddress)
				files, err := client.GetFiles("integration_test")
				Expect(err).To(BeNil())
				Expect(files).To(ContainSubstring(golook.GolookSystem.UUID))

				client.Host = fmt.Sprintf("http://%s:8383", d1.Container.NetworkSettings.IPAddress)
				files, err = client.GetFiles("integration_test")
				Expect(err).To(BeNil())
				Expect(files).To(ContainSubstring(golook.GolookSystem.UUID))
			})
		})

		It("informs ONE peers about files in folders", func() {
			int.RunPeerInDocker(func(client *docker.Client, container *docker.Container) {

				routerFiles.NewPeer(routing.NewKey(container.NetworkSettings.IPAddress), container.NetworkSettings.IPAddress)
				routerSystem.NewPeer(routing.NewKey(container.NetworkSettings.IPAddress), container.NetworkSettings.IPAddress)

				systemService.Run()
				_, err := fileServices.Report(&models.FileReport{Path: "integration_test.go"})
				if err != nil {
					log.WithError(err).Fatal("Cannot make file report in test.")
				}

				result, _ := fileServices.Query("integration_test")

				log.WithField("test", "integration").Debug("Files retrieved:"+
					"num=%d", len(result.(fileQueryData)))

				var files = result.(fileQueryData)

				Expect(files).To(HaveKey(golook.GolookSystem.UUID))
			})
		})
	})

	Context("broadcast report service", func() {
		It("informs peers about files in folders", func() {
			int.RunPeerInDocker(func(client *docker.Client, container *docker.Container) {

				routerFiles.NewPeer(routing.NewKey(container.NetworkSettings.IPAddress), container.NetworkSettings.IPAddress)
				routerSystem.NewPeer(routing.NewKey(container.NetworkSettings.IPAddress), container.NetworkSettings.IPAddress)

				systemService.Run()

				result, err := fileServices.Report(&models.FileReport{Path: "integration_test.go"})

				Expect(err).To(BeNil())
				expectedKey, err := filepath.Abs("integration_test.go")
				if err != nil {
					log.WithField("file", "integration_test.go").Fatal("Cannot generate absolute filepath.")
				}
				Expect(result).To(HaveKey(expectedKey))
			})
		})
	})

	Context("system service", func() {

		It("ensures that all peers know about all other peers.", func() {
			int.RunPeersInDocker(2, func(dockerizedPeers []*int.DockerizedGolook) {

				for i := range dockerizedPeers {
					container := dockerizedPeers[i].Container
					routerFiles.NewPeer(routing.NewKey(container.NetworkSettings.IPAddress), container.NetworkSettings.IPAddress)
					routerSystem.NewPeer(routing.NewKey(container.NetworkSettings.IPAddress), container.NetworkSettings.IPAddress)
				}

				//second time: simulates the initial call of Run...
				systemService.Run()
				//second time: simulates the next scheduled call of Run...
				systemService.Run()

				time.Sleep(time.Second)

				log.WithField("test", "integration_systems").Error(utils.MarshalSD(GetSystems()))

				l, _ := client.GetLog()
				log.Info(l)

				// query one peer at random if ALL systems are registered with the neighbor
				container := dockerizedPeers[0].Container
				client.Host = fmt.Sprintf("http://%s:8383", container.NetworkSettings.IPAddress)
				system, _ := client.GetSystem()
				log.WithField("test", "integration_systems").Info(system)

				Expect(system).To(ContainSubstring(golook.GolookSystem.UUID))
			})
		})

		It("can inform peers about the system.", func() {
			int.RunPeerInDocker(func(c *docker.Client, container *docker.Container) {
				routerFiles.NewPeer(routing.NewKey(container.NetworkSettings.IPAddress), container.NetworkSettings.IPAddress)
				routerSystem.NewPeer(routing.NewKey(container.NetworkSettings.IPAddress), container.NetworkSettings.IPAddress)

				result := systemService.broadcastSystem(
					peerSystemReport{
						Uuid: golook.GolookSystem.UUID,
						System: map[string]*golook.System{
							golook.GolookSystem.UUID: golook.GolookSystem,
						},
						IsDeletion: false,
					},
				)

				// query the peer if the system is registered with the neighbor
				client.Host = fmt.Sprintf("http://%s:8383", container.NetworkSettings.IPAddress)
				system, _ := client.GetSystem()
				log.WithField("test", "integration").Debug(system)

				Expect(result).ToNot(BeNil())
				Expect(system).To(ContainSubstring(golook.GolookSystem.UUID))
			})
		})
	})
})
