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
)

type testRouteLayerClient struct {
	message string
}

func (t *testRouteLayerClient) Handle(method string, m interface{}) interface{} {
	t.message = m.(string)
	return method
}

var _ = Describe("The route layer callback registrar", func() {

	It("buffer added callbacks", func() {
		RouteLayerRegistrar.RegisterDefaultClient(&testRouteLayerClient{})
		Expect(RouteLayerRegistrar.callback).ToNot(BeNil())
	})

	It("calls the route layer when tasked to do so", func() {
		t := &testRouteLayerClient{}
		RouteLayerRegistrar.RegisterDefaultClient(t)

		res := ToRouteLayer("test", "msg")

		Expect(t.message).To(Equal("msg"))
		Expect(res.(string)).To(Equal("test"))
	})

	It("rejects messages if the route layer is not active", func() {
		RouteLayerRegistrar.RegisterDefaultClient(nil)

		res := ToRouteLayer("test", "msg")

		Expect(res).To(BeNil())
	})
})
