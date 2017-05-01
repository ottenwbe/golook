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
	"github.com/ottenwbe/golook/broker/runtime"
	"github.com/satori/go.uuid"
)

var _ = Describe("The duplicate filter", func() {
	It("detects dulicates from the same system.", func() {
		source1 := Source{1, runtime.GolookSystem.UUID}
		source2 := Source{1, runtime.GolookSystem.UUID}

		Expect(duplicateMap.CheckForDuplicates(source1)).To(BeFalse())
		Expect(duplicateMap.CheckForDuplicates(source2)).To(BeTrue())
	})
	It("ignores duplicated ids from different systems", func() {
		source1 := Source{1, uuid.NewV4().String()}
		source2 := Source{1, uuid.NewV4().String()}

		Expect(duplicateMap.CheckForDuplicates(source1)).To(BeFalse())
		Expect(duplicateMap.CheckForDuplicates(source2)).To(BeFalse())
	})
})
