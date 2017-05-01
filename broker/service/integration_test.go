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

	"github.com/fsouza/go-dockerclient"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	log "github.com/sirupsen/logrus"

	"github.com/ottenwbe/golook/broker/models"
	"github.com/ottenwbe/golook/broker/routing"
	"github.com/ottenwbe/golook/broker/runtime"
	"github.com/ottenwbe/golook/broker/utils"
	int "github.com/ottenwbe/golook/test/integration"
)

var _ = Describe("The service layer's", func() {

	var (
		systemService *SystemService
	)

	BeforeEach(func() {
		systemService = &SystemService{}
	})

	Context("broadcast query service", func() {
		It("informs peers about files in folders", func() {
			int.RunPeerInDocker(func(client *docker.Client, container *docker.Container) {
				reportService := newReportService(BCastReport)
				queryService := newQueryService(BCastQueries)
				broadCastRouter.NewPeer(routing.NewKey(container.NetworkSettings.IPAddress), container.NetworkSettings.IPAddress)
				systemService.Run()
				_, err := reportService.MakeFileReport(&models.FileReport{Path: "integration_test.go"})
				if err != nil {
					log.Fatal("Cannot make file report in test.")
				}

				result, _ := queryService.MakeFileQuery("integration_test")

				var files map[string][]*models.File
				err = utils.Unmarshal(result, &files)
				if err != nil {
					log.Fatal("Cannot unmarshal files in test.")
				}

				Expect(files).To(HaveKey(runtime.GolookSystem.UUID))
			})
		})
	})

	Context("broadcast report service", func() {
		It("informs peers about files in folders", func() {
			int.RunPeerInDocker(func(client *docker.Client, container *docker.Container) {
				reportService := newReportService(BCastReport)
				broadCastRouter.NewPeer(routing.NewKey(container.NetworkSettings.IPAddress), container.NetworkSettings.IPAddress)
				systemService.Run()

				result, err := reportService.MakeFileReport(&models.FileReport{Path: "integration_test.go"})

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
				broadCastRouter.NewPeer(routing.NewKey(container.NetworkSettings.IPAddress), container.NetworkSettings.IPAddress)

				result := systemService.broadcastSystem(false)

				Expect(result).ToNot(BeNil())
				var peerResponse PeerResponse
				err := utils.Unmarshal(result, &peerResponse)
				Expect(err).To(BeNil())
				Expect(peerResponse.Error).To(BeFalse())
			})
		})
	})
})