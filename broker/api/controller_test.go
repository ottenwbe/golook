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

package api

import (
	"net/http"
	"net/http/httptest"

	"github.com/ottenwbe/golook/broker/models"
	. "github.com/ottenwbe/golook/broker/runtime"

	"bytes"
	"encoding/json"
	"github.com/gorilla/mux"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/ottenwbe/golook/broker/service"
	"sync"
)

var _ = Describe("The management endpoint", func() {

	var (
		rr *httptest.ResponseRecorder
		m  *mux.Router

		testHTTPCall = func(req *http.Request, path string, f func(http.ResponseWriter, *http.Request)) {
			rr = httptest.NewRecorder()
			m = mux.NewRouter()
			m.HandleFunc(path, f)
			m.ServeHTTP(rr, req)
		}
	)

	BeforeEach(func() {
		service.CloseFileServices()
		service.OpenFileServices(service.MockFileServices)
	})

	AfterEach(func() {
		service.CloseFileServices()
		service.OpenFileServices(service.BroadcastFiles)
	})

	Context(HTTPApiEndpoint, func() {
		It("should return the api information", func() {
			ApplyConfiguration() //Ensure server is up

			req, err := http.NewRequest(http.MethodGet, HTTPApiEndpoint, nil)
			testHTTPCall(req, HTTPApiEndpoint, getAPI)

			Expect(err).To(BeNil())
			Expect(rr.Code).To(Equal(http.StatusOK))
		})

		It("should return an error when the server is not running."+
			"However, this should not happen, since in this case the endpoint should not be called at all.", func() {
			HTTPServer = nil
			defer ApplyConfiguration()

			req, _ := http.NewRequest(http.MethodGet, HTTPApiEndpoint, nil)
			testHTTPCall(req, HTTPApiEndpoint, getAPI)

			Expect(rr.Code).To(Equal(http.StatusInternalServerError))
		})

		It("should return an error when the http server's information cannot be marshalled.", func() {
			HTTPServer = &testHTTPServer{}
			defer ApplyConfiguration()

			req, _ := http.NewRequest(http.MethodGet, HTTPApiEndpoint, nil)
			testHTTPCall(req, HTTPApiEndpoint, getAPI)

			Expect(rr.Code).To(Equal(http.StatusInternalServerError))
		})
	})

	Context(LogEndpoint, func() {
		It("should return no errors, when ....", func() {
			req, err := http.NewRequest(http.MethodGet, LogEndpoint, nil)
			testHTTPCall(req, LogEndpoint, getLog)

			Expect(err).To(BeNil())
			Expect(rr.Code).To(Equal(http.StatusOK))
		})
	})

	Context(InfoEndpoint, func() {

		It("should return the current app info", func() {
			req, err := http.NewRequest(http.MethodPut, InfoEndpoint, nil)
			testHTTPCall(req, InfoEndpoint, getInfo)

			Expect(err).To(BeNil())
			Expect(rr.Code).To(Equal(http.StatusOK))
			Expect(rr.Body.String()).To(Equal(EncodeAppInfo(NewAppInfo())))
		})
	})

	Context(FileEndpoint, func() {

		It("should return an error 400, when the request's body is empty", func() {
			req, err := http.NewRequest(http.MethodPut, FileEndpoint, nil)
			testHTTPCall(req, FileEndpoint, putFile)

			Expect(err).To(BeNil())
			Expect(rr.Code).To(Equal(http.StatusBadRequest))
		})

		It("should return an error 400, when the body is invalid", func() {
			b := []byte("{invalid}")

			req, err := http.NewRequest(http.MethodPut, FileEndpoint, bytes.NewBuffer(b))
			testHTTPCall(req, FileEndpoint, putFile)

			Expect(err).To(BeNil())
			Expect(service.ReportService.(*service.MockReportService).FileReport).To(BeNil())
			Expect(rr.Code).To(Equal(http.StatusBadRequest))
		})

		It("should call the report service to report the file", func() {

			fp := &models.FileReport{
				Path: "test.txt",
			}

			b, _ := json.Marshal(fp)
			req, err := http.NewRequest(http.MethodPut, FileEndpoint, bytes.NewBuffer(b))
			testHTTPCall(req, FileEndpoint, putFile)

			Expect(err).To(BeNil())
			Expect(rr.Code).To(Equal(http.StatusOK))
			Expect(*service.ReportService.(*service.MockReportService).FileReport).To(Equal(*fp))

		})

		It("should call the query service", func() {

			req, err := http.NewRequest(http.MethodPut, FileEndpoint+"/test.txt", nil)
			testHTTPCall(req, QueryEndpoint, getFiles)

			Expect(err).To(BeNil())
			Expect(rr.Code).To(Equal(http.StatusOK))
			Expect(rr.Body.String()).To(Equal("{}"))
			Expect(service.QueryService.(*service.MockQueryService).SearchString).To(Equal("test.txt"))
		})
	})
})

type testHTTPServer struct {
}

func (*testHTTPServer) StartServer(_ *sync.WaitGroup) {
}

func (*testHTTPServer) Info() map[string]interface{} {
	return map[string]interface{}{"{": make(chan bool)}
}
