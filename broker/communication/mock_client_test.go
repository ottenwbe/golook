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

var _ = Describe("The mock client", func() {

	It("provides a builder function.", func() {
		m := newMockClient()
		Expect(m).ToNot(BeNil())
		Expect(reflect.TypeOf(m).String()).To(Equal(reflect.TypeOf(&MockClient{}).String()))
	})

	It("should record the name of the latest call() and the number of calls.", func() {
		m := &MockClient{}
		m.Call("test", nil)
		Expect(m.Name).To(Equal("test"))
		Expect(m.VisitedCall).To(Equal(1))
	})

	It("should record the number of calls to Url().", func() {
		m := &MockClient{}
		m.Url()
		m.Url()
		Expect(m.VisitedUrl).To(Equal(2))
	})

})
