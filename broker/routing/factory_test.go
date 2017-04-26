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
		r := NewRouter("test")

		Expect(r).ToNot(BeNil())
		Expect(reflect.TypeOf(r)).To(Equal(reflect.TypeOf(&FloodingRouter{})))
	})

	It("registers and deregisters a created router", func() {
		r := NewRouter("test")
		Expect(r).ToNot(BeNil())
		Expect(communication.MessageDispatcher.HasHandler("test")).To(BeTrue())
		DeleteRouter("test")
		Expect(communication.MessageDispatcher.HasHandler("test")).To(BeFalse())
	})
})
