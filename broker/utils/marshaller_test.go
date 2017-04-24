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
package utils

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Marshalling a message", func() {

	It(" should be the same result after demarshalling when using the bytes method", func() {
		testString := "test"

		msg, err1 := MarshalB(testString)
		var s string
		err2 := UnmarshalB(msg, &s)

		Expect(err1).To(BeNil())
		Expect(err2).To(BeNil())
		Expect(s).To(Equal(testString))
	})

	It(" should be the same result after demarshalling when using the string method", func() {
		testString := "test"

		msg, err1 := MarshalS(testString)
		var s string
		err2 := UnmarshalS(msg, &s)

		Expect(err1).To(BeNil())
		Expect(err2).To(BeNil())
		Expect(s).To(Equal(testString))
	})

	It("should return an error when channels or other unsupported types are used", func() {
		c := make(chan bool)

		_, err := MarshalB(c)

		Expect(err).ToNot(BeNil())
	})

	//TODO own context
	It("should return an error when unmarshalling a faulty input", func() {
		s := ""
		type test struct {
			i int
		}
		var t test

		err := UnmarshalB([]byte(s), &t)

		Expect(err).ToNot(BeNil())
	})
})
