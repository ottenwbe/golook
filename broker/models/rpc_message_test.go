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
package models

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Marshalling a message", func() {
	It("it should be the same result after demarshalling", func() {
		testString := "test"

		msg, err1 := Marshal(testString)
		var s string
		err2 := Unmarshal(msg, &s)

		Expect(err1).To(BeNil())
		Expect(err2).To(BeNil())
		Expect(s).To(Equal(testString))
	})

	It("should throw an error when channels or other unsupported types are marshalled", func() {
		c := make(chan bool)

		_, err := Marshal(c)

		Expect(err).ToNot(BeNil())
	})
})

var _ = Describe("The encapsulated message", func() {
	It("should comprise a method name and the content after its creation", func() {
		m, err := NewRpcMessage("method", "msg")
		Expect(err).To(BeNil())
		Expect(m.Method).To(Equal("method"))
		Expect(len(m.Content)).ToNot(BeZero())
	})

	It("should support to get the encapsulated method", func() {
		m, err := NewRpcMessage("method", "msg")

		var s string
		m.GetEncapsulated(&s)

		Expect(err).To(BeNil())
		Expect(s).To(Equal("msg"))
	})
})
