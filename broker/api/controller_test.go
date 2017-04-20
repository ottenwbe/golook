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
	"github.com/gorilla/mux"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"net/http"
	"net/http/httptest"
)

var _ = Describe("The management endpoint", func() {

	var (
		rr *httptest.ResponseRecorder
		m  *mux.Router
	)

	testHttpCall := func(req *http.Request, path string, f func(http.ResponseWriter, *http.Request)) {
		rr = httptest.NewRecorder()
		m = mux.NewRouter()
		m.HandleFunc(path, f)
		m.ServeHTTP(rr, req)
	}

	Context("/file", func() {
		It("should return a 400 and nack, when body is empty", func() {
			req, err := http.NewRequest(http.MethodPut, "/file", nil)
			testHttpCall(req, "/file", putFile)

			Expect(err).To(BeNil())
			Expect(rr.Code).To(Equal(http.StatusBadRequest))
			Expect(rr.Body.String()).To(ContainSubstring(Nack))
		})
	})

	Context("/folder", func() {
		It("should return a 400 and nack, when body is empty", func() {
			req, err := http.NewRequest(http.MethodPut, "/folder", nil)
			testHttpCall(req, "/folder", putFolder)

			Expect(err).To(BeNil())
			Expect(rr.Code).To(Equal(http.StatusBadRequest))
			Expect(rr.Body.String()).To(ContainSubstring(Nack))
		})
	})
})
