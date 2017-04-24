////Copyright 2016-2017 Beate Ottenwälder
////
////Licensed under the Apache License, Version 2.0 (the "License");
////you may not use this file except in compliance with the License.
////You may obtain a copy of the License at
////
////http://www.apache.org/licenses/LICENSE-2.0
////
////Unless required by applicable law or agreed to in writing, software
////distributed under the License is distributed on an "AS IS" BASIS,
////WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
////See the License for the specific language governing permissions and
////limitations under the License.
package routing

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"reflect"
)

var _ = Describe("The mocked router", func() {

	var (
		router *MockedLookRouter
	)

	BeforeEach(func() {
		router = NewMockedRouter().(*MockedLookRouter)
	})

	Context("logic of mock router", func() {
		It("should increase the Visited flag with each call to Handle", func() {
			router.Handle("test_handle", nil)
			Expect(router.Visited).To(Equal(1))
			Expect(router.VisitedMethod).To(Equal("test_handle"))
		})

		It("should increase the Visited flag with each call to Route", func() {
			router.Route(SysKey(), "test_route", nil)
			Expect(router.Visited).To(Equal(1))
			Expect(router.VisitedMethod).To(Equal("test_route"))
		})

		It("should increase the Visited flag with each call to HandlerFunction", func() {
			router.HandlerFunction("test_handler", nil)
			Expect(router.Visited).To(Equal(1))
			Expect(router.VisitedMethod).To(Equal("test_handler"))
		})

		It("should increase the Visited flag with each call to BroadCast", func() {
			router.BroadCast("test_bc", nil)
			Expect(router.Visited).To(Equal(1))
			Expect(router.VisitedMethod).To(Equal("test_bc"))
		})

		It("should return 'mock' as name", func() {
			Expect(router.Name()).To(Equal("mock"))
		})
	})

	Context("RunWithMockedRouter", func() {
		It("should replace a router with a mock, execute a function in a block with the replaced router,"+
			" and then reset the original router afterwards", func() {
			index := NewRouter("r")
			RunWithMockedRouter(&index, func() {
				Expect(reflect.TypeOf(index)).ToNot(Equal(reflect.TypeOf(NewRouter("r"))))
				//Test if original router was successfully set
				Expect(reflect.TypeOf(index)).To(Equal(reflect.TypeOf(NewMockedRouter())))
			})
			//Test if reset to original router was successful
			Expect(reflect.TypeOf(index)).To(Equal(reflect.TypeOf(NewRouter("r"))))
		})
	})
})
