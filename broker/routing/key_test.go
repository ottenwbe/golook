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
	"github.com/satori/go.uuid"
)

var _ = Describe("The (routing) key", func() {
	It("should have a unique nil value.", func() {
		key := NilKey()
		Expect(key.id).To(Equal(uuid.Nil))
	})

	It("should create unique keys based on a name.", func() {
		key1 := NewKey("test")
		key2 := NewKey("test")
		Expect(key1).To(Equal(key2))
	})

	It("should create unique keys based on a name.", func() {
		key1 := NewKeyN(uuid.Nil, "test")
		key2 := NewKeyN(uuid.Nil, "test")
		Expect(key1).To(Equal(key2))
	})

	It("should not alter the uuid when created with a uuid.", func() {
		testUUID := uuid.NewV4()
		key1 := NewKeyU(testUUID)
		Expect(key1.id).To(Equal(testUUID))
	})
})
