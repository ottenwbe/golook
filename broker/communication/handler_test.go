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
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/ottenwbe/golook/broker/models"
)

var _ = Describe("The router callback registrar", func() {

	It("stores routers that have been registered", func() {
		const TEST_NAME = "test1"
		MessageDispatcher.RegisterHandler(TEST_NAME, &testRouteLayerClient{}, testMsg{}, testResponse{})
		Expect(*MessageDispatcher).ToNot(BeNil())
		Expect((*MessageDispatcher)[TEST_NAME]).ToNot(BeNil())
	})

	It("allows to query for stored routers", func() {
		const TEST_NAME = "testQuery"
		MessageDispatcher.RegisterHandler(TEST_NAME, &testRouteLayerClient{}, testMsg{}, testResponse{})
		Expect(MessageDispatcher.HasHandler(TEST_NAME)).To(BeTrue())
	})

	It("calls the router when tasked to do so", func() {
		const (
			MSG_CONTENT = "msg"
			TEST_NAME   = "test"
		)
		t := &testRouteLayerClient{}
		MessageDispatcher.RegisterHandler(TEST_NAME, t, testMsg{}, testResponse{})

		res, err := MessageDispatcher.handleMessage(TEST_NAME, &testMsgConv{MSG_CONTENT})

		Expect(err).To(BeNil())
		Expect(t.message).To(Equal(MSG_CONTENT))
		Expect(res.(string)).To(Equal(TEST_NAME))
	})

	It("rejects messages, i.e., returns nil, if a handler is not registered", func() {
		const TEST_NAME = "should_not_exist_test"
		MessageDispatcher.RegisterHandler(TEST_NAME, nil, testMsg{}, testResponse{})

		res, err := MessageDispatcher.handleMessage(TEST_NAME, &testMsgConv{"msg"})

		Expect(err).ToNot(BeNil())
		Expect(res).To(BeNil())
	})

	It("rejects messages, i.e., returns nil, if a handler is removed befor a message is sent.", func() {
		const (
			msgContent = "msg"
			testName   = "test"
		)
		t := &testRouteLayerClient{}
		MessageDispatcher.RegisterHandler(testName, t, testMsg{}, testResponse{})

		MessageDispatcher.RemoveHandler(testName)
		res, err := MessageDispatcher.handleMessage(testName, &testMsgConv{msgContent})

		Expect(err).ToNot(BeNil())
		Expect(res).To(BeNil())
	})
})

type (
	testRouteLayerClient struct {
		message string
	}
	testMsg struct {
	}
	testResponse struct {
	}
	testMsgConv struct {
		testString string
	}
)

func (tConverter *testMsgConv) Unmarshal(v interface{}) error {
	*v.(*string) = tConverter.testString
	return nil
}

func (t *testRouteLayerClient) Handle(method string, params models.EncapsulatedValues) interface{} {
	params.Unmarshal(&t.message)
	return method
}
