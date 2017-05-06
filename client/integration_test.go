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

package client

import (
	"fmt"
	"github.com/fsouza/go-dockerclient"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/ottenwbe/golook/broker/models"
	int "github.com/ottenwbe/golook/test/integration"
)

var _ = Describe("The api endpoint ", func() {

	Context(httpAPIEndpoint, func() {
		It("can be called when a golook instance is running.", func() {
			int.RunPeerInDocker(
				func(client *docker.Client, container *docker.Container) {
					Host = fmt.Sprintf("http://%s:8383", container.NetworkSettings.IPAddress)

					apiInfo, err := GetAPI()

					Expect(err).To(BeNil())
					Expect(apiInfo).To(ContainSubstring("endpoints"))
				})
		})
	})

	Context(logEndpoint, func() {
		It("retrieves the log of a running instance.", func() {
			int.RunPeerInDocker(
				func(client *docker.Client, container *docker.Container) {
					Host = fmt.Sprintf("http://%s:8383", container.NetworkSettings.IPAddress)

					logInfo, err := GetLog()

					Expect(err).To(BeNil())
					Expect(logInfo).To(ContainSubstring("Starting"))
				})
		})
	})

	Context(infoEndpoint, func() {
		It("returns the system report.", func() {
			int.RunPeerInDocker(
				func(client *docker.Client, container *docker.Container) {
					Host = fmt.Sprintf("http://%s:8383", container.NetworkSettings.IPAddress)

					info, err := GetInfo()

					Expect(err).To(BeNil())
					Expect(info).To(ContainSubstring("name"))
				})
		})
	})

	Context(fileEndpoint, func() {
		It("accepts file reports.", func() {
			int.RunPeerInDocker(
				func(client *docker.Client, container *docker.Container) {
					Host = fmt.Sprintf("http://%s:8383", container.NetworkSettings.IPAddress)

					fileReport := models.FileReport{Path: ".", Delete: false}
					configInfo, err := ReportFiles(fileReport)

					Expect(err).To(BeNil())
					Expect(configInfo).To(ContainSubstring("name"))
				})
		})
	})

	Context(systemEndpoint, func() {
		It("retrieves the stored systems of a running instance.", func() {
			int.RunPeerInDocker(
				func(client *docker.Client, container *docker.Container) {
					Host = fmt.Sprintf("http://%s:8383", container.NetworkSettings.IPAddress)

					systems, err := GetSystem()

					Expect(err).To(BeNil())
					Expect(systems).To(ContainSubstring("name"))
				})
		})
	})

	Context(configEndpoint, func() {
		It("retrieves the config of a running instance.", func() {
			int.RunPeerInDocker(
				func(client *docker.Client, container *docker.Container) {
					Host = fmt.Sprintf("http://%s:8383", container.NetworkSettings.IPAddress)

					configInfo, err := GetConfig()

					Expect(err).To(BeNil())
					Expect(configInfo).To(ContainSubstring("log.level"))
				})
		})
	})

})
