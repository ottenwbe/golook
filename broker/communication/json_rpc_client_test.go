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

package communication

// Tests based on https://github.com/ybbus/jsonrpc/blob/master/jsonrpc_test.go

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/sirupsen/logrus"
	"github.com/ybbus/jsonrpc"
)

type TestParams struct {
	A int    `json:"a"`
	B string `json:"b"`
}

var _ = Describe("The rpc client", func() {

	const (
		EXPECTED_RESPONSE_CONTENT = "1-2-3"
	)

	var (
		requestChan     chan string
		httpServer      *httptest.Server
		testClient      RpcClient
		expectedContent []byte
	)

	BeforeEach(func() {

		var err error
		if expectedContent, err = json.Marshal(TestParams{A: 1, B: "test"}); err != nil {
			logrus.WithError(err).Fatal("Test failed.")
		}

		requestChan = make(chan string, 1)

		// start the test server
		httpServer = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			data, _ := ioutil.ReadAll(r.Body)
			defer r.Body.Close()
			// put request and body to channel for the client to investigate them
			requestChan <- string(data)

			b, _ := json.Marshal(EXPECTED_RESPONSE_CONTENT)
			fmt.Fprintf(w, `{"jsonrpc":"2.0","result":%s,"id":0}`, string(b))
		}))

		testClient = &JsonRpcClientStub{
			serverUrl: httpServer.URL,
			c:         jsonrpc.NewJsonRPCClient(fmt.Sprintf("%s/rpc", httpServer.URL)),
		}
	})

	AfterEach(func() {
		httpServer.Close()
	})

	It("should send a message to the server and in turn receive the response without an error", func() {
		var expectedString string

		// make a test call to the server (see httpServer)
		result, err := testClient.Call("testMethod", TestParams{A: 1, B: "test"})
		Expect(err).To(BeNil())

		// unmarshal the result
		result.Unmarshal(&expectedString)

		// get the msg which has been received by the server
		res := <-requestChan
		Expect(res).To(ContainSubstring(string(expectedContent)))
		Expect(expectedString).To(Equal(EXPECTED_RESPONSE_CONTENT))
	})

	It("should return an error when an invalid type is to be transferred as content, e.g. a channel", func() {
		_, err := testClient.Call("no method", make(chan bool))
		Expect(err).ToNot(BeNil())
	})
})
