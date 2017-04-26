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
	"encoding/json"
	"github.com/fsouza/go-dockerclient"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/ottenwbe/golook/broker/api"
	"github.com/ottenwbe/golook/broker/communication"
	"github.com/ottenwbe/golook/broker/models"
	"github.com/ottenwbe/golook/broker/routing"
	"github.com/ottenwbe/golook/broker/runtime"
	"github.com/ottenwbe/golook/broker/service"
	"github.com/ottenwbe/golook/broker/utils"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
)

var _ = Describe("The golook server's", func() {

	Context("http api", func() {
		It("allows clients to query the application's /info endpoint from a server running in a docker container", func() {
			RunInDockerContainer(func(client *docker.Client, container *docker.Container) {

				appInfo := &runtime.AppInfo{}
				compareInfo := runtime.NewAppInfo()
				containerInfo := GetContainerInfo(client, container)

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
	})
})

var _ = Describe("The golook server's", func() {
	Context("rpc api", func() {
		It("registers a system via the rest api", func() {
			RunInDockerContainer(func(client *docker.Client, container *docker.Container) {
				var (
					containerInfo = GetContainerInfo(client, container)
					testRouter    = routing.NewRouter("system")
					testRpcClient = communication.NewRPCClient("http://" + containerInfo.NetworkSettings.IPAddress + ":8382")
				)
				testRouter.NewPeer(routing.NewKey(containerInfo.NetworkSettings.MacAddress), testRpcClient)

				testSystem := &models.SystemFiles{System: runtime.GolookSystem, Files: nil}
				res := testRouter.BroadCast(service.SYSTEM_REPORT,
					&service.SystemTransfer{Uuid: runtime.GolookSystem.UUID, System: testSystem, IsDeletion: false})

				actualResponse := service.PeerResponse{}
				err := utils.Unmarshal(res, &actualResponse)

				Expect(err).To(BeZero())
				Expect(actualResponse.Error).To(BeFalse())
			})
		})

		It("deletes a registered system via the rest api", func() {
			RunInDockerContainer(func(client *docker.Client, container *docker.Container) {
				var (
					containerInfo = GetContainerInfo(client, container)
					testRouter    = routing.NewRouter("system")
					testRpcClient = communication.NewRPCClient("http://" + containerInfo.NetworkSettings.IPAddress + ":8382")
				)
				testRouter.NewPeer(routing.NewKey(containerInfo.NetworkSettings.MacAddress), testRpcClient)

				testSystem := &models.SystemFiles{System: runtime.GolookSystem, Files: nil}
				testRouter.BroadCast(service.SYSTEM_REPORT,
					&service.SystemTransfer{Uuid: runtime.GolookSystem.UUID, System: testSystem, IsDeletion: false})
				res := testRouter.BroadCast(service.SYSTEM_REPORT,
					&service.SystemTransfer{Uuid: runtime.GolookSystem.UUID, System: testSystem, IsDeletion: true})

				actualResponse := service.PeerResponse{}
				err := utils.Unmarshal(res, &actualResponse)

				Expect(err).To(BeZero())
				Expect(actualResponse.Error).To(BeFalse())
			})
		})

		It("handles and stores a file in order to query it afterwards", func() {
			RunInDockerContainer(func(client *docker.Client, container *docker.Container) {
				var (
					containerInfo = GetContainerInfo(client, container)
					testRouter    = routing.NewRouter("system")
					testRpcClient = communication.NewRPCClient("http://" + containerInfo.NetworkSettings.IPAddress + ":8382")
				)
				testRouter.NewPeer(routing.NewKey(containerInfo.NetworkSettings.MacAddress), testRpcClient)

				testSystem := &models.SystemFiles{System: runtime.GolookSystem, Files: nil}
				testRouter.BroadCast(service.SYSTEM_REPORT,
					&service.SystemTransfer{Uuid: runtime.GolookSystem.UUID, System: testSystem, IsDeletion: false})
				defer testRouter.BroadCast(service.SYSTEM_REPORT,
					&service.SystemTransfer{Uuid: runtime.GolookSystem.UUID, System: testSystem, IsDeletion: true})

				f, err := models.NewFile("integration_test.go")
				Expect(err).To(BeZero())

				res_fr := testRouter.BroadCast(service.FILE_REPORT,
					&service.PeerFileReport{
						Files:   map[string]*models.File{f.ShortName: f},
						Replace: false,
						System:  runtime.GolookSystem.UUID,
					})

				actualResponse := service.PeerResponse{}
				err = utils.Unmarshal(res_fr, &actualResponse)

				Expect(err).To(BeZero())
				Expect(actualResponse.Error).To(BeFalse())

				res_q := testRouter.BroadCast(service.FILE_QUERY, service.PeerFileQuery{
					SearchString: "integration_test.go",
				})

				actualResponse2 := service.PeerResponse{}
				err = utils.Unmarshal(res_q, &actualResponse2)

				logrus.Info("Result of query: " + res_q.(string))

				var files map[string]*models.SystemFiles
				err = utils.Unmarshal(actualResponse2.Data, &files)

				Expect(err).To(BeZero())
				Expect(actualResponse2.Error).To(BeFalse())

				_, ok := files[runtime.GolookSystem.UUID]
				Expect(ok).To(BeTrue())

				_, ok2 := files[runtime.GolookSystem.UUID].Files[f.Name]
				Expect(ok2).To(BeTrue())
			})
		})

	})
})
