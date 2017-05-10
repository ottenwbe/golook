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
	"reflect"
)

var _ = Describe("The communication factory", func() {

	AfterEach(func() {
		ClientType = jsonRPC
	})

	It("can build json rpc clients; when 'jsonRPC' is configured.", func() {
		ClientType = jsonRPC
		client := NewRPCClient("1.2.3.4")
		Expect(reflect.TypeOf(client).String()).To(Equal(reflect.TypeOf(&JsonRpcClientStub{}).String()))
	})

	It("can build mock rpc clients, when 'MockRPC' is configured.", func() {
		ClientType = MockRPC
		client := NewRPCClient("1.2.3.4")
		Expect(reflect.TypeOf(client).String()).To(Equal(reflect.TypeOf(&MockClient{}).String()))
	})

})
