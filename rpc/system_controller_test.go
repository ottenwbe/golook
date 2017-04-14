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
	"github.com/gorilla/mux"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/ottenwbe/golook/global"

	"net/http"
	"net/http/httptest"
)

var _ = Describe("The management endpoint", func() {

	var (
		rr *httptest.ResponseRecorder
		m  *mux.Router
	)

	f := func(req *http.Request, path string, f func(http.ResponseWriter, *http.Request)) {
		rr = httptest.NewRecorder()
		m = mux.NewRouter()
		m.HandleFunc(path, f)
		m.ServeHTTP(rr, req)
	}

	Context(EP_INFO, func() {
		It("should return the app info", func() {
			req, err := http.NewRequest("GET", EP_INFO, nil)
			f(req, EP_INFO, getInfo)

			Expect(err).To(BeNil())
			Expect(rr.Code).To(Equal(http.StatusOK))
			Expect(rr.Body.String()).To(ContainSubstring(VERSION))
		})
	})

	Context(EP_SYSTEM, func() {
		It("should return the system info with GET", func() {
			req, err := http.NewRequest("GET", EP_SYSTEM, nil)
			f(req, EP_SYSTEM, getSystem)

			Expect(err).To(BeNil())
			Expect(rr.Code).To(Equal(http.StatusOK))
			Expect(rr.Body.String()).To(ContainSubstring("os"))
		})
	})
})
