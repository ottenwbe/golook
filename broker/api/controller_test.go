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
	"bytes"
	"encoding/json"
	"github.com/gorilla/mux"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/ottenwbe/golook/broker/models"
	golook "github.com/ottenwbe/golook/broker/runtime/core"
	"github.com/ottenwbe/golook/broker/service"
	"github.com/ottenwbe/golook/utils"
	"net/http"
	"net/http/httptest"
	"sync"
)

var _ = Describe("The management endpoint", func() {

	var (
		rr *httptest.ResponseRecorder
		m  *mux.Router

		testHTTPCall = func(req *http.Request, path string, f func(http.ResponseWriter, *http.Request)) *httptest.ResponseRecorder {
			rr = httptest.NewRecorder()
			m = mux.NewRouter()
			m.HandleFunc(path, f)
			m.ServeHTTP(rr, req)
			return rr
		}
	)

	BeforeEach(func() {
		service.CloseFileServices(fileServices)
		fileServices = service.OpenFileServices(service.MockFileServices)
	})

	AfterEach(func() {
		service.CloseFileServices(fileServices)
		fileServices = service.OpenFileServices(service.BroadcastFiles)
	})

	Context(SystemEndpoint, func() {
		It("should return the system's stored on this machine", func() {
			req, err := http.NewRequest(http.MethodGet, SystemEndpoint, nil)
			testHTTPCall(req, SystemEndpoint, getSystem)

			Expect(err).To(BeNil())
			Expect(rr.Code).To(Equal(http.StatusOK))
		})
	})

	Context(ConfigEndpoint, func() {
		It("should return the configuration", func() {

			req, err := http.NewRequest(http.MethodGet, ConfigEndpoint, nil)
			testHTTPCall(req, ConfigEndpoint, getConfiguration)

			Expect(err).To(BeNil())
			Expect(rr.Code).To(Equal(http.StatusOK))

		})

		It("should return a status code 500, when the configuration comprises a value that cannot be"+
			"marshalled to json, i.e., the return value", func() {
			var testService service.ConfigurationService = &testInvalidConfigService{}
			utils.Mock(&configurationService, &testService, func() {
				req, err := http.NewRequest(http.MethodGet, ConfigEndpoint, nil)
				testHTTPCall(req, ConfigEndpoint, getConfiguration)

				Expect(err).To(BeNil())
				Expect(rr.Code).To(Equal(http.StatusInternalServerError))
			})
		})

		It("should return a status code 500, when the configuration service is not set.", func() {

			tmp := configurationService
			configurationService = nil
			defer func() {
				configurationService = tmp
			}()

			req, err := http.NewRequest(http.MethodGet, ConfigEndpoint, nil)
			testHTTPCall(req, ConfigEndpoint, getConfiguration)

			Expect(err).To(BeNil())
			Expect(rr.Code).To(Equal(http.StatusInternalServerError))
		})
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
		It("should return a status 200 when the log can be fetched.", func() {
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

			result := &golook.AppInfo{}
			err = utils.Unmarshal(rr.Body.String(), &result)
			if err != nil {
				Fail("Cannot marshal reslut")
			}

			Expect(err).To(BeNil())
			Expect(rr.Code).To(Equal(http.StatusOK))
			Expect(*result).To(Equal(*golook.NewAppInfo()))
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
			Expect(service.AccessMockedReportService(fileServices).FileReport).To(BeNil())
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
			Expect(*service.AccessMockedReportService(fileServices).FileReport).To(Equal(*fp))

		})

		It("should call the query service", func() {

			req, err := http.NewRequest(http.MethodPut, FileEndpoint+"/test.txt", nil)
			testHTTPCall(req, QueryEndpoint, getFiles)

			Expect(err).To(BeNil())
			Expect(rr.Code).To(Equal(http.StatusOK))
			Expect(rr.Body.String()).To(Equal("{}"))
			Expect(service.AccessMockedQueryService(fileServices).SearchString).To(Equal("test.txt"))
		})
	})
})

type testHTTPServer struct {
}

func (*testHTTPServer) Stop() {
}

func (*testHTTPServer) Name() string {
	return ""
}

func (*testHTTPServer) IsRunning() bool {
	return true
}

func (*testHTTPServer) StartServer(_ *sync.WaitGroup) {
}

func (*testHTTPServer) Info() map[string]interface{} {
	return map[string]interface{}{"{": make(chan bool)}
}

type testInvalidConfigService struct {
}

func (*testInvalidConfigService) GetConfiguration() map[string]map[string]interface{} {
	return map[string]map[string]interface{}{
		"a": {"b": make(chan bool, 0)},
	}
}
