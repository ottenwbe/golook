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

import (
	"encoding/json"

	"github.com/ottenwbe/golook/broker/models"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/sirupsen/logrus"
)

type TestServerParams struct {
	A int    `json:"a"`
	B string `json:"b"`
}

type testHandler struct {
	msg []TestServerParams
}

func (t *testHandler) Handle(method string, params models.MsgParams) interface{} {
	logrus.Infof("Rec params: %s", params.(*RpcReceiverParams).params)
	params.GetObject(&t.msg)
	return TestServerParams{}

}

var _ = Describe("The rpc server", func() {

	It("should have an rpc handler that take a method and call a corresponding registered handler", func() {
		testData := []TestServerParams{{A: 1, B: "2"}}
		rpcHandler := &RpcHandler{"test"}
		correspondingHandler := &testHandler{}

		RouterRegistrar.RegisterClient("test", correspondingHandler, TestServerParams{}, TestServerParams{})

		params, err := json.Marshal(testData)
		m := json.RawMessage(params)
		_, errJ := rpcHandler.ServeJSONRPC(nil, &m)

		Expect(err).To(BeNil())
		Expect(errJ).To(BeNil())
		Expect(correspondingHandler.msg).To(Equal(testData))
	})
})
