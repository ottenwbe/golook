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
		RouterRegistrar.RegisterClient(TEST_NAME, &testRouteLayerClient{}, testMsg{}, testResponse{})
		Expect(RouterRegistrar.routerCallback).ToNot(BeNil())
		Expect(RouterRegistrar.routerCallback[TEST_NAME]).ToNot(BeNil())
	})

	It("allows to query for stored routers", func() {
		const TEST_NAME = "testQuery"
		RouterRegistrar.RegisterClient(TEST_NAME, &testRouteLayerClient{}, testMsg{}, testResponse{})
		Expect(RouterRegistrar.HasClient(TEST_NAME)).To(BeTrue())
	})

	It("calls the router when tasked to do so", func() {
		const (
			MSG_CONTENT = "msg"
			TEST_NAME   = "test"
		)
		t := &testRouteLayerClient{}
		RouterRegistrar.RegisterClient(TEST_NAME, t, testMsg{}, testResponse{})

		res := toRouteLayer(TEST_NAME, &testMsgConv{MSG_CONTENT})

		Expect(t.message).To(Equal(MSG_CONTENT))
		Expect(res.(string)).To(Equal(TEST_NAME))
	})

	It("rejects messages, i.e., returns nil, if the router is not registered", func() {
		const TEST_NAME = "atest"
		RouterRegistrar.RegisterClient(TEST_NAME, nil, testMsg{}, testResponse{})

		res := toRouteLayer(TEST_NAME, &testMsgConv{"msg"})

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

func (tConverter *testMsgConv) GetObject(v interface{}) error {
	*v.(*string) = tConverter.testString
	return nil
}

func (t *testRouteLayerClient) Handle(method string, params models.MsgParams) interface{} {
	params.GetObject(&t.message)
	return method
}
