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

var _ = Describe("The mock template", func() {

	It("should allow to run methods in a block with a mocked value for an interface", func() {
		var test int64 = 5
		var i int64 = 8
		Mock(&i, &test, func() {
			Expect(i).To(Equal(int64(5)))
		})
		Expect(i).To(Equal(int64(8)))
	})

	It("should recover after panicing", func() {
		var test int64 = 5
		var i int64 = 8
		Mock(&i, &test, func() {
			panic("controlled panic of a test")
		})
		Expect(i).To(Equal(int64(8)))
	})
})
