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
		const testName = "test1"
		MessageDispatcher.RegisterHandler(testName, &testRouteLayerClient{}, testMsg{}, testResponse{})
		Expect(*MessageDispatcher).ToNot(BeNil())
		Expect((*MessageDispatcher)[testName]).ToNot(BeNil())
	})

	It("allows to query for stored routers", func() {
		const testName = "testQuery"
		MessageDispatcher.RegisterHandler(testName, &testRouteLayerClient{}, testMsg{}, testResponse{})
		Expect(MessageDispatcher.HasHandler(testName)).To(BeTrue())
	})

	It("calls the router when tasked to do so", func() {
		const (
			msgContent = "msg"
			testName   = "test"
		)
		t := &testRouteLayerClient{}
		MessageDispatcher.RegisterHandler(testName, t, testMsg{}, testResponse{})

		res, err := MessageDispatcher.handleMessage(testName, &testMsgConv{msgContent})

		Expect(err).To(BeNil())
		Expect(t.message).To(Equal(msgContent))
		Expect(res.(string)).To(Equal(testName))
	})

	It("rejects messages, i.e., returns nil, if a handler is not registered", func() {
		const testName = "should_not_exist_test"
		MessageDispatcher.RegisterHandler(testName, nil, testMsg{}, testResponse{})

		res, err := MessageDispatcher.handleMessage(testName, &testMsgConv{"msg"})

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
