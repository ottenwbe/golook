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
	"errors"
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/ottenwbe/golook/broker/models"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
)

var _ = Describe("The client for the api endpoint ", func() {

	const expectedDefaultResult = `{"test":"1.2","data":dat,"id":0}`

	var (
		// start the test servers
		httpServer = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ioutil.ReadAll(r.Body)
			defer r.Body.Close()

			fmt.Fprint(w, expectedDefaultResult)
		}))

		httpErrorServer = httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, r *http.Request) {
			ioutil.ReadAll(r.Body)
			defer r.Body.Close()

			http.Error(writer, errors.New("test error").Error(), http.StatusNotImplemented)
		}))
	)

	Context(httpAPIEndpoint, func() {
		It("returns the result received from the server.", func() {
			Host = httpServer.URL
			response, err := GetAPI()
			Expect(err).To(BeNil())
			Expect(response).To(Equal(expectedDefaultResult))
		})

		It("returns an error, when the server returns a http error.", func() {
			Host = httpErrorServer.URL
			_, err := GetAPI()
			Expect(err).ToNot(BeNil())
		})

		It("returns an error, when the host's url is not set.", func() {
			Host = "http://127.0.0"
			_, error := GetAPI()
			Expect(error).ToNot(BeNil())
		})
	})

	Context(fileEndpoint, func() {
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
