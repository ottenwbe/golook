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
)

var _ = Describe("The management endpoint", func() {

	var (
		rr *httptest.ResponseRecorder
		m  *mux.Router

		testHttpCall = func(req *http.Request, path string, f func(http.ResponseWriter, *http.Request)) {
			rr = httptest.NewRecorder()
			m = mux.NewRouter()
			m.HandleFunc(path, f)
			m.ServeHTTP(rr, req)
		}

		savedReportService = reportService
		savedQueryService  = queryService
	)

	BeforeEach(func() {
		reportService = &testService{}
		queryService = &testService{}
	})

	AfterEach(func() {
		reportService = savedReportService
		queryService = savedQueryService
	})

	Context(INFO_EP, func() {

		It("should return the current app info", func() {
			req, err := http.NewRequest(http.MethodPut, INFO_EP, nil)
			testHttpCall(req, INFO_EP, getInfo)

			Expect(err).To(BeNil())
			Expect(rr.Code).To(Equal(http.StatusOK))
			Expect(rr.Body.String()).To(Equal(EncodeAppInfo(NewAppInfo())))
		})
	})

	Context(FILE_EP, func() {

		It("should return an error 400 and nack, when body is empty", func() {
			req, err := http.NewRequest(http.MethodPut, FILE_EP, nil)
			testHttpCall(req, FILE_EP, putFile)

			Expect(err).To(BeNil())
			Expect(rr.Code).To(Equal(http.StatusBadRequest))
			Expect(rr.Body.String()).To(ContainSubstring(NACK))
		})

		It("should return an error 400 and nack, when the body is invalid", func() {
			b := []byte("{invalid}")

			req, err := http.NewRequest(http.MethodPut, FILE_EP, bytes.NewBuffer(b))
			testHttpCall(req, FILE_EP, putFile)

			Expect(err).To(BeNil())
			Expect(reportService.(*testService).fileReport).To(BeNil())
			Expect(rr.Code).To(Equal(http.StatusBadRequest))
			Expect(rr.Body.String()).To(ContainSubstring(NACK))
		})

		It("should call the report service to report the file", func() {

			fp := &models.FileReport{
				Path:    "test.txt",
				Monitor: false,
				Replace: false,
			}

			b, _ := json.Marshal(fp)
			req, err := http.NewRequest(http.MethodPut, FILE_EP, bytes.NewBuffer(b))
			testHttpCall(req, FILE_EP, putFile)

			Expect(err).To(BeNil())
			Expect(rr.Code).To(Equal(http.StatusOK))
			Expect(rr.Body.String()).To(Equal(ACK))
			Expect(*reportService.(*testService).fileReport).To(Equal(*fp))

		})

		It("should call the query service", func() {

			req, err := http.NewRequest(http.MethodPut, FILE_EP+"/test.txt", nil)
			testHttpCall(req, FILE_QUERY_EP, getFiles)

			Expect(err).To(BeNil())
			Expect(rr.Code).To(Equal(http.StatusOK))
			Expect(rr.Body.String()).To(Equal("{}"))
			Expect(queryService.(*testService).searchString).To(Equal("test.txt"))
		})
	})

	Context(FOLDER_EP, func() {
		It("should return an error 400 and nack, when body is empty", func() {
			req, err := http.NewRequest(http.MethodPut, FOLDER_EP, nil)
			testHttpCall(req, FOLDER_EP, putFolder)

			Expect(err).To(BeNil())
			Expect(rr.Code).To(Equal(http.StatusBadRequest))
			Expect(rr.Body.String()).To(ContainSubstring(NACK))
		})

		It("should call the report service to report the folder", func() {

			testFileReport := &models.FileReport{
				Path:    "a_folder",
				Monitor: false,
				Replace: false,
			}
			expectedFileReport := *testFileReport

			b, _ := json.Marshal(testFileReport)
			req, err := http.NewRequest(http.MethodPut, FOLDER_EP, bytes.NewBuffer(b))
			testHttpCall(req, FOLDER_EP, putFolder)

			Expect(err).To(BeNil())
			Expect(rr.Code).To(Equal(http.StatusOK))
			Expect(rr.Body.String()).To(Equal(ACK))
			Expect(*reportService.(*testService).folderReport).To(Equal(expectedFileReport))

		})

	})
})

type testService struct {
	searchString string
	fileReport   *models.FileReport
	folderReport *models.FileReport
}

func (ts *testService) MakeFileQuery(searchString string) interface{} {
	ts.searchString = searchString
	return "{}"
}

func (ts *testService) MakeFileReport(fileReport *models.FileReport) {
	ts.fileReport = fileReport
}

func (ts *testService) MakeFolderReport(folderReport *models.FileReport) {
	ts.folderReport = folderReport
}
