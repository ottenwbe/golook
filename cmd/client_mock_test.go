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
package cmd

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe(" MockGoLookClient ", func() {
	It("should be set when a function is run in context of 'RunWithMockedGolookClient'", func() {
		RunWithMockedGolookClient(func() {
			Expect(testConvertMockGolookClient).ToNot(Panic())
		})
	})

	It("should be set when a function is run in context of 'RunWithMockedGolookClientF'", func() {
		RunWithMockedGolookClientF(func() {
			Expect(testConvertMockGolookClient).ToNot(Panic())
			Expect(GolookClient.(*MockGolookClient).fileName).To(Equal("fileName.txt"))
			Expect(GolookClient.(*MockGolookClient).folderName).To(Equal("folder"))
		}, "fileName.txt", "folder")
	})
})

func testConvertMockGolookClient() {
	_ = GolookClient.(*MockGolookClient)
}
