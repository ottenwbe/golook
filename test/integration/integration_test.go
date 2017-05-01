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
	"encoding/json"
	"github.com/fsouza/go-dockerclient"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/ottenwbe/golook/broker/api"
	"github.com/ottenwbe/golook/broker/runtime"
	"io/ioutil"
	"net/http"
)

var _ = Describe("The golook server's", func() {

	Context("http api", func() {
		It("allows clients to query the application's /info endpoint from a server running in a docker container", func() {
			RunPeerInDocker(func(client *docker.Client, container *docker.Container) {

				appInfo := &runtime.AppInfo{}
				compareInfo := runtime.NewAppInfo()
				//containerInfo := GetContainerInfo(client, container)

				// make test request
				r, errGet := http.Get("http://" + container.NetworkSettings.IPAddress + ":8383" + api.InfoEndpoint)
				Expect(errGet).To(BeNil())

				// get the result
				b, errRead := ioutil.ReadAll(r.Body)
				Expect(errRead).To(BeNil())
				errMarshal := json.Unmarshal(b, appInfo)
				Expect(errMarshal).To(BeNil())

				// verify that the result is correct
				Expect(appInfo).ToNot(BeNil())
				Expect(appInfo.Version).To(Equal(compareInfo.Version))
				Expect(appInfo.System.IP).To(Equal(container.NetworkSettings.IPAddress))
			})
		})
	})
})
