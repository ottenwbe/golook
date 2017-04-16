package rpc_server

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/gorilla/mux"

	. "github.com/ottenwbe/golook/app"

	"net/http/httptest"
	"net/http"
)

var _ = Describe("The file_management endpoint", func() {

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
})
