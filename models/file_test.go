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
	"encoding/json"
	"io"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	log "github.com/sirupsen/logrus"
)

const (
	testFileName = "test.txt"
)

var _ = Describe("The model", func() {

	var (
		fileChannel chan *File
		pipeReader  *io.PipeReader
		pipeWriter  *io.PipeWriter
	)

	BeforeEach(func() {
		fileChannel = make(chan *File)
		pipeReader, pipeWriter = io.Pipe()
	})

	Context("File", func() {
		It("can be marshalled and unmarshalled, however, meta information is ignored", func() {
			go PipeWriteFile(pipeWriter, newTestFile(testFileName))
			go PipeReadFile(pipeReader, fileChannel)

			result := <-fileChannel

			Expect(result).ToNot(BeNil())
			Expect(result.Name).To(Equal(expectedResultFile(testFileName).Name))
			Expect(result.Meta).To(Equal(expectedResultFile(testFileName).Meta))
		})
	})

	Context("Files", func() {
		It("can be marshalled and unmarshalled, however, meta information is ignored", func() {
			go PipeWriteFiles(pipeWriter, map[string]File{testFileName: newTestFile(testFileName), testFileName + "2": newTestFile(testFileName + "2")})
			go PipeReadFiles(pipeReader, fileChannel)

			result := <-fileChannel

			Expect(result).ToNot(BeNil())
			Expect(result.Name).To(Equal(expectedResultFile(testFileName).Name))
			Expect(result.Meta).To(Equal(expectedResultFile(testFileName).Meta))
		})
	})
})

func PipeWriteFiles(pipeWriter *io.PipeWriter, files map[string]File) {
	b, _ := json.Marshal(files)
	defer pipeWriter.Close()
	pipeWriter.Write(b)
}

func PipeReadFiles(pipeReader *io.PipeReader, c chan *File) {
	if f, err := DecodeFiles(pipeReader); err == nil && f[testFileName].Name == testFileName {
		tmpF := f[testFileName]
		c <- &tmpF
	} else {
		log.Printf("Error expected nil, got %s", err)
		log.Printf("File, expected %s got %s", testFileName, f[testFileName].Name)
	}
	c <- nil
}

func PipeWriteFile(pipeWriter *io.PipeWriter, file File) {
	b, _ := json.Marshal(file)
	defer pipeWriter.Close()
	pipeWriter.Write(b)
}

func PipeReadFile(pipeReader *io.PipeReader, c chan *File) {
	if f, err := DecodeFile(pipeReader); err == nil && f.Name == testFileName && f.Meta.Monitor == false {
		c <- &f
	} else {
		log.Printf("Error expected nil, got %s", err)
		log.Printf("File, expected %s got %s", testFileName, f.Name)
	}
	c <- nil
}

func newTestFile(fileName string) File {
	f := File{}
	f.Name = fileName
	f.Accessed = time.Time{}
	f.Created = time.Time{}
	f.Modified = time.Time{}
	f.Meta.Monitor = true
	return f
}

func expectedResultFile(fileName string) File {
	f := newTestFile(fileName)
	f.Meta.Monitor = false
	return f
}
