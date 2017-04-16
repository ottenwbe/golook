package report

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/gorilla/mux"

	. "github.com/ottenwbe/golook/rpc_server"

	"net/http/httptest"
	"net/http"
	"testing"
)

func TestReports(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Report Controller Suite")
}

var _ = Describe("The file_management endpoint", func() {

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

	Context(EP_REPORT+"/file", func() {
		It("should return a 400 and nack, when body is empty", func() {
			req, err := http.NewRequest(http.MethodPut, EP_REPORT, nil)
			testHttpCall(req, EP_REPORT, putFile)

			Expect(err).To(BeNil())
			Expect(rr.Code).To(Equal(http.StatusBadRequest))
			Expect(rr.Body.String()).To(ContainSubstring(Nack))
		})
	})

	Context(EP_REPORT+"/folder", func() {
		It("should return a 400 and nack, when body is empty", func() {
			req, err := http.NewRequest(http.MethodPut, EP_REPORT, nil)
			testHttpCall(req, EP_REPORT, putFolder)

			Expect(err).To(BeNil())
			Expect(rr.Code).To(Equal(http.StatusBadRequest))
			Expect(rr.Body.String()).To(ContainSubstring(Nack))
		})
	})
})
