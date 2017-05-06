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
	"path/filepath"

	"fmt"
	"github.com/fsouza/go-dockerclient"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/ottenwbe/golook/broker/models"
	"github.com/ottenwbe/golook/broker/routing"
	"github.com/ottenwbe/golook/broker/runtime"
	"github.com/ottenwbe/golook/client"
	int "github.com/ottenwbe/golook/test/integration"
	log "github.com/sirupsen/logrus"
	"time"
)

var _ = Describe("The service layer's", func() {

	var (
		systemService *SystemService
		fileServices  fileServices
		routerSystem  *router
		routerFiles   *router
	)

	BeforeEach(func() {
		fileServices = newFileServices(BroadcastFiles)
		fileServices.open()
		routerFiles = fileServices.(*scenarioBroadcastFiles).r

		systemService = newSystemService()
		routerSystem = systemService.router
	})

	AfterEach(func() {
		systemService.close()
		fileServices.close()
	})

	Context("broadcast query service", func() {
		It("informs ALL peers about files in folders", func() {
			d1 := &int.DockerizedGolook{}
			d2 := &int.DockerizedGolook{}
			d1.Init()
			defer d1.Stop()
			d2.Init()
			defer d2.Stop()
			d1.Start()
			d2.Start()

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

			// query one of the neighbor hosts
			client.Host = fmt.Sprintf("http://%s:8383", d1.Container.NetworkSettings.IPAddress)
			files, err := client.GetFiles("integration_test")
			log.Info(files)

			// query one of the neighbor hosts
			client.Host = fmt.Sprintf("http://%s:8383", d1.Container.NetworkSettings.IPAddress)
			system, err := client.GetSystem()
			log.Info(system)

			Expect(err).To(BeNil())
			Expect(files).To(ContainSubstring(runtime.GolookSystem.UUID))
		})
	})

	Context("broadcast query service", func() {
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
					"num=%d", len(result.(FileQueryData)))

				var files = result.(FileQueryData)

				Expect(files).To(HaveKey(runtime.GolookSystem.UUID))
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
		It("can inform peers about the system.", func() {
			int.RunPeerInDocker(func(client *docker.Client, container *docker.Container) {
				routerFiles.NewPeer(routing.NewKey(container.NetworkSettings.IPAddress), container.NetworkSettings.IPAddress)
				routerSystem.NewPeer(routing.NewKey(container.NetworkSettings.IPAddress), container.NetworkSettings.IPAddress)

				log.Info("Sending broadcast")
				result := systemService.broadcastSystem(
					PeerSystemReport{
						Uuid: runtime.GolookSystem.UUID,
						System: map[string]*runtime.System{
							runtime.GolookSystem.UUID: runtime.GolookSystem,
						},
						IsDeletion: false,
					},
				)

				Expect(result).ToNot(BeNil())
			})
		})
	})
})
