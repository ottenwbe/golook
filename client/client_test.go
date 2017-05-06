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
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/ottenwbe/golook/broker/models"
)

var _ = Describe("The api endpoint ", func() {

	Context(HTTPApiEndpoint, func() {
		It("returns an error, when the host's url is not set.", func() {
			Host = "http://127.0.0"
			_, error := GetApi()
			Expect(error).ToNot(BeNil())
		})
	})

	Context(FileEndpoint, func() {
		It("returns an error, when the host's url is not set correctly.", func() {
			Host = "http://127.0.0"
			_, error := ReportFiles(models.FileReport{Path: "test", Delete: false})
			Expect(error).ToNot(BeNil())
		})
		It("returns an error, when the host's url is not set at all.", func() {
			Host = ""
			_, error := ReportFiles(models.FileReport{Path: "test", Delete: false})
			Expect(error).ToNot(BeNil())
		})
	})
})
