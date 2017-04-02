package utils

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"path/filepath"
)

var _ = Describe("The file's", func() {

	const FILE_NAME = "file_io_test.go"

	It(",metadata can be read", func() {
		fp, errFilePath := filepath.Abs(FILE_NAME)

		file, errFile := NewFile(fp)

		Expect(errFilePath).To(BeNil())
		Expect(errFile).To(BeNil())
		Expect(file).To(Not(BeNil()))
		Expect(file.Name).To(Equal(FILE_NAME))
	})

})
