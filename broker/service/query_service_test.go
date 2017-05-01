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

package service

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/ottenwbe/golook/broker/routing"
	"reflect"
)

var _ = Describe("The query service", func() {
	It("creates a local query service by default", func() {
		q := newQueryService(LocalQueries)
		Expect(q).ToNot(BeNil())
		Expect(reflect.TypeOf(q)).To(Equal(reflect.TypeOf(&localQueryService{})))
	})

	It("creates broadcast query service when BCastQueries is specified", func() {
		q := newQueryService(BCastQueries)
		Expect(q).ToNot(BeNil())
		Expect(reflect.TypeOf(q)).To(Equal(reflect.TypeOf(&broadcastQueryService{})))
	})

	It("calls the routing service to initiate the query", func() {
		tmp := routing.Router(broadCastRouter)
		routing.RunWithMockedRouter(&tmp, func() {
			q := broadcastQueryService{}
			q.MakeFileQuery("test.txt")

			Expect(routing.AccessMockedRouter(broadCastRouter).Visited).To(Equal(1))
			Expect(routing.AccessMockedRouter(broadCastRouter).VisitedMethod).To(Equal(FILE_QUERY))
		})
	})

})
