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
package file_management

import (
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"path/filepath"
)

var _ = Describe("The file's", func() {

	const FILE_NAME = "file_io_test.go"

	It("metadata of file file_io_test.go can be read", func() {
		fp, errFilePath := filepath.Abs(FILE_NAME)

		file, errFile := NewFile(fp)

		Expect(errFilePath).To(BeNil())
		Expect(errFile).To(BeNil())
		Expect(file).To(Not(BeNil()))
		Expect(file.Name).To(Equal(FILE_NAME))
	})

	It("non existing file is not read", func() {
		fp := fmt.Sprintf("%s_does_not_exist", FILE_NAME)

		_, errFile := NewFile(fp)

		Expect(errFile).ToNot(BeNil())
	})

})
