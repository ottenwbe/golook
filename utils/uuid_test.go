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

var _ = Describe("uuids", func() {

	var (
		uuid1, uuid2 string
		err1, err2   error
	)

	BeforeEach(func() {
		uuid1, err1 = NewUUID()
		uuid2, err2 = NewUUID()
	})

	It("generated at random should differ", func() {
		Expect(err1).To(BeNil())
		Expect(err2).To(BeNil())
		Expect(uuid1).To(Not(Equal(uuid2)))
	})
})
