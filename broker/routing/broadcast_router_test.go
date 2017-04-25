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
	"github.com/ottenwbe/golook/broker/models"
	"reflect"
)

var _ = Describe("The broadcast router", func() {
	It("implements the 'Router' interface", func() {
		r := newBroadcastRouter("test")

		Expect(r).ToNot(BeNil())
		Expect(reflect.TypeOf(r)).To(Equal(reflect.TypeOf(&BroadcastRouter{})))
	})

	It("broadcasts messages to one or more peers", func() {
		testPeers := []*mockPeer{{}, {}}

		r := newBroadcastRouter("test")
		r.NewPeer(NewKey("peer1"), testPeers[0])
		r.NewPeer(NewKey("peer2"), testPeers[1])

		r.BroadCast("test", 123)

		for i := range testPeers {
			Expect(testPeers[i].visitedCall).To(Equal(1))
			Expect(testPeers[i].request.Method).To(Equal("test"))
			Expect(testPeers[i].request.Params).To(Equal("123"))
		}
	})

	It("should broadcast instead of sending directed messages via 'Route'", func() {
		testPeers := []*mockPeer{{}, {}}

		r := newBroadcastRouter("test")
		r.NewPeer(NewKey("peer1"), testPeers[0])
		r.NewPeer(NewKey("peer2"), testPeers[1])

		r.Route(NewKey("peer2"), "test", 123)

		for i := range testPeers {
			Expect(testPeers[i].visitedCall).To(Equal(1))
			Expect(testPeers[i].request.Method).To(Equal("test"))
			Expect(testPeers[i].request.Params).To(Equal("123"))
		}
	})
})

type mockPeer struct {
	visitedCall int
	request     *RequestMessage
}

func (p *mockPeer) Call(index string, message interface{}) (models.MsgParams, error) {
	p.visitedCall += 1
	p.request = message.(*RequestMessage)
	return nil, nil
}

func (mockPeer) Url() string {
	return "test"
}
