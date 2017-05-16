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

package routing

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/ottenwbe/golook/broker/communication"
	"reflect"
)

var _ = Describe("The router factory", func() {
	It("creates a new router and registeres it", func() {
		r := NewRouter("test", BroadcastRouterType)

		Expect(r).ToNot(BeNil())
		Expect(reflect.TypeOf(r)).To(Equal(reflect.TypeOf(&BroadCastRouter{})))
	})

	It("registers a created router with the communication layer", func() {
		r := NewRouter("test", BroadcastRouterType)

		ActivateRouter(r)
		defer DeactivateRouter(r)

		Expect(communication.MessageDispatcher.HasHandler("test")).To(BeTrue())
	})

	It("adds default peers.", func() {
		const expectedPeer = "1.1.1.1"
		defaultPeers = []string{expectedPeer}
		r := NewRouter("test", MockRouterType)
		Expect(AccessMockedRouter(r).PeerName).To(Equal(expectedPeer))
		defaultPeers = []string{}
	})

	It("deregisters am registered router from the communication layer", func() {
		r := NewRouter("test", BroadcastRouterType)
		ActivateRouter(r)

		DeactivateRouter(r)

		Expect(communication.MessageDispatcher.HasHandler("test")).ToNot(BeTrue())
	})

	It("returns nil if the router type is not known.", func() {
		r := NewRouter("test", "testRouter")
		Expect(r).To(BeNil())
	})

})
