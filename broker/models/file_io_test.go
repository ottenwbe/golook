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
	"fmt"
	"path/filepath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	log "github.com/sirupsen/logrus"
)

var _ = Describe("The file's", func() {

	const fileName = "file_io_test.go"

	It("metadata of file file_io_test.go can be read", func() {
		fp, errFilePath := filepath.Abs(fileName)
		if errFilePath != nil {
			log.WithError(errFilePath).Fatal("File I/O test failed")
		}

		file, errFile := NewFile(fp)

		Expect(errFile).To(BeNil())
		Expect(file).To(Not(BeNil()))
		Expect(file.Name).To(Equal(fp))
		Expect(file.Meta.State).To(Equal(Created))
	})

	It("non existing file is not read and marked as removed", func() {
		fp := fmt.Sprintf("%s_does_not_exist", fileName)

		file, errFile := NewFile(fp)

		Expect(errFile).To(BeNil())
		Expect(file.Meta.State).To(Equal(Removed))
	})

})
