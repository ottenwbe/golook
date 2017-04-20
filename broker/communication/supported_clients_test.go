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

var _ = Describe("SupportedClients", func() {

	BeforeEach(func() {
		LookupClients = newCommunicationClients()
	})

	AfterEach(func() {
		LookupClients = newCommunicationClients()
	})

	testClientBuilder := func(url string) LookupClient { return NewLookupRPCClient(url) }

	It("should be able to store a builder and use this builder to create the LookClients by a key", func() {
		key := "rpc"
		LookupClients.Add(key, testClientBuilder, true)
		builtClient, ok := LookupClients.Build("no_url", key)

		Expect(ok).To(BeTrue())
		Expect(builtClient).ToNot(BeNil())
	})

	It("should panic when a duplicated entry is about to be made", func() {
		key := "rpc"
		LookupClients.Add(key, testClientBuilder, true)
		Expect(func() {
			LookupClients.Add(key, testClientBuilder, false)
		}).To(Panic())
	})

	It("should panic when two default clients are selected", func() {
		key1 := "rpc_1"
		key2 := "rpc_2"
		LookupClients.Add(key1, testClientBuilder, true)
		Expect(func() {
			LookupClients.Add(key2, testClientBuilder, true)
		}).To(Panic())
	})

	It("should return all selected communication clients", func() {
		key := "rpc"
		LookupClients.Add(key, testClientBuilder, true)
		supported := SupportedCommunicationClients()
		Expect(len(supported)).To(Equal(1))
		Expect(supported).To(Equal([]string{key}))
	})

})
