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
)

var _ = Describe("The rpc params", func() {

	It("can be unmarshalled", func() {
		expectedResult := []aType{{123}}

		b, err := json.Marshal(expectedResult)
		if err != nil {
			Fail("Error while preparig a test")
		}
		jsonRPCParams := JsonRPCParams{
			params: json.RawMessage(b),
		}

		// test unmarshalling the data
		var testResult aType
		err = jsonRPCParams.Unmarshal(&testResult)

		Expect(err).To(BeNil())
		Expect(testResult).To(Equal(expectedResult[0]))
	})

	It("returns an error when a slice should be unmarshalled", func() {
		expectedResult := []aType{{123}, {234}}

		b, err := json.Marshal(expectedResult)
		if err != nil {
			Fail("Error while preparig a test")
		}
		jsonRPCParams := JsonRPCParams{
			params: json.RawMessage(b),
		}

		// test unmarshalling the data
		var testResult aType
		err = jsonRPCParams.Unmarshal(&testResult)

		Expect(err).ToNot(BeNil())
	})
})

var _ = Describe("The rpc server", func() {

	It("can dispatch a message to a corresponding handler, when said handler exists", func() {
		testMessage := []TestServerParams{{A: 1, B: "2"}}
		rpcServer := &JsonRPCServerStub{"test", true}
		correspondingHandler := &testHandler{}

		MessageDispatcher.RegisterHandler("test", correspondingHandler, TestServerParams{}, TestServerParams{})

		params, err := json.Marshal(testMessage)
		m := json.RawMessage(params)
		_, errJ := rpcServer.ServeJSONRPC(nil, &m)

		Expect(err).To(BeNil())
		Expect(errJ).To(BeNil())
		Expect(correspondingHandler.msg).To(Equal(testMessage[0]))
	})

	It("retuns an error when the method is finalized.", func() {
		testMessage := []TestServerParams{{A: 1, B: "2"}}
		rpcServer := &JsonRPCServerStub{"test", true}
		correspondingHandler := &testHandler{}

		MessageDispatcher.RegisterHandler("test", correspondingHandler, TestServerParams{}, TestServerParams{})

		params, err := json.Marshal(testMessage)
		Expect(err).To(BeNil())

		rpcServer.Finalize()

		m := json.RawMessage(params)
		_, errJ := rpcServer.ServeJSONRPC(nil, &m)
		Expect(errJ).ToNot(BeNil())
	})

	It("associates handler functions that are called.", func() {
		testMessage := []TestServerParams{{A: 1, B: "2"}}
		rpcServer := &JsonRPCServerStub{"test", true}
		correspondingHandler := &testHandler{}

		MessageDispatcher.RegisterHandler("test", correspondingHandler, TestServerParams{}, TestServerParams{})

		params, err := json.Marshal(testMessage)
		Expect(err).To(BeNil())

		rpcServer.Finalize()

		m := json.RawMessage(params)
		_, errJ := rpcServer.ServeJSONRPC(nil, &m)
		Expect(errJ).ToNot(BeNil())
	})
})

type TestServerParams struct {
	A int    `json:"a"`
	B string `json:"b"`
}

type testHandler struct {
	msg TestServerParams
}

func (t *testHandler) Handle(method string, params models.EncapsulatedValues) interface{} {
	params.Unmarshal(&t.msg)
	return TestServerParams{}

}

type aType struct {
	A int `json:"a"`
}
