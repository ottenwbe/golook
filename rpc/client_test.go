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
package rpc

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/gorilla/mux"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/ottenwbe/golook/global"
	"github.com/ottenwbe/golook/utils"
)

const (
	sysName   = "system"
	FILE_NAME = "file.txt"
)

var _ = Describe("The rpc", func() {

	var (
		server *httptest.Server
		client *LookupClientData
	)

	BeforeEach(func() {
		tmpClient := NewLookClient("http://127.0.0.1", 8123)
		client = tmpClient.(*LookupClientData)
	})

	AfterEach(func() {
		// ensure that the close method is executed and not forgotten
		server.Close()
		client = nil
	})

	Context(" System Methods ", func() {
		It("should return a valid system with Get", func() {
			server = httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, _ *http.Request) {
				s := newTestSystem()
				bytes, _ := json.Marshal(s)
				fmt.Fprintln(writer, string(bytes))
			}))
			client.serverUrl = server.URL

			result, err := client.DoGetSystem(sysName)
			Expect(err).To(BeNil())
			Expect(result).To(Not(BeNil()))
			Expect(result.Name).To(Equal(sysName))
		})

		It("should return a nil system with Get when the server does not exist", func() {
			client.serverUrl = "/"
			result, err := client.DoGetSystem(sysName)
			Expect(result).To(BeNil())
			Expect(err).ToNot(BeNil())
		})

		It("should send a valid system to the server with Put", func() {

			testSystem := newTestSystem()

			server = httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {
				receivedSystem, _ := DecodeSystem(req.Body)
				Expect(receivedSystem.Name).To(Equal(testSystem.Name))
			}))
			client.serverUrl = server.URL

			result := client.DoPutSystem(testSystem)

			Expect(result).To(Not(BeNil()))
		})

		It("should transfer the delete request for a specific system to the server with DELETE", func() {

			testSystemName := "testSystem"

			server = httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {
				params := mux.Vars(req)
				system := params["system"]
				Expect(system).To(Equal(testSystemName))
			}))
			server.URL = "/systems/{system}"
			client.serverUrl = server.URL
			client.systemName = testSystemName

			result := client.DoDeleteSystem()

			Expect(result).To(Not(BeNil()))
		})
	})

	Context(" File Methods ", func() {
		It("should return a valid set of files with Get", func() {
			server := httptest.NewServer(
				http.HandlerFunc(
					func(writer http.ResponseWriter, _ *http.Request) {
						b, _ := json.Marshal(map[string]utils.File{FILE_NAME: newTestFile()})
						fmt.Fprint(writer, string(b))
					}))
			client.serverUrl = server.URL

			files, err := client.DoGetFiles("testSystem")
			Expect(err).To(BeNil())
			Expect(len(files)).To(Equal(1))
			Expect(files[FILE_NAME].Name).To(Equal(FILE_NAME))
		})
	})

	Context("Get Home", func() {

		const testString = "TestString"

		It("should pass the string which was sent by a server to the calle of DoGetHome()", func() {

			server := httptest.NewServer(
				http.HandlerFunc(
					func(writer http.ResponseWriter, _ *http.Request) {
						fmt.Fprintln(writer, testString)
					}))
			client.serverUrl = server.URL

			Expect(client.DoGetHome()).To(Equal(testString + "\n"))
		})
	})

})

func newTestSystem() *System {
	return &System{
		Name:  sysName,
		Files: nil,
		IP:    "1.1.1.1",
		OS:    "linux",
		UUID:  "1234"}
}

func newTestFile() utils.File {
	return utils.File{
		Name: FILE_NAME,
	}
}
