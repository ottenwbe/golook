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
		chanbools  chan bool
		pipeReader *io.PipeReader
		pipeWriter *io.PipeWriter
	)

	BeforeEach(func() {
		chanbools = make(chan bool)
		pipeReader, pipeWriter = io.Pipe()
	})

	Context("File", func() {
		It("can be marshalled and unmarshalled", func() {
			go PipeWriteFile(pipeWriter)
			go PipeReadFile(pipeReader, chanbools)
			result := <-chanbools
			Expect(result).To(BeTrue())
		})
	})

	Context("Files", func() {
		It("can be marshalled and unmarshalled", func() {
			go PipeWriteFiles(pipeWriter)
			go PipeReadFiles(pipeReader, chanbools)
			result := <-chanbools
			Expect(result).To(BeTrue())
		})
	})
})

func PipeWriteFiles(pipeWriter *io.PipeWriter) {
	b, _ := json.Marshal(map[string]File{testFileName: newTestFile(testFileName), testFileName + "2": newTestFile(testFileName)})
	defer pipeWriter.Close()
	pipeWriter.Write(b)
}

func PipeReadFiles(pipeReader *io.PipeReader, c chan bool) {
	if f, err := DecodeFiles(pipeReader); err == nil && f[testFileName].Name == testFileName {
		c <- true
	} else {
		log.Printf("Error expected nil, got %s", err)
		log.Printf("File, expected %s got %s", testFileName, f[testFileName].Name)
	}
	c <- false
}

func PipeWriteFile(pipeWriter *io.PipeWriter) {
	b, _ := json.Marshal(newTestFile(testFileName))
	defer pipeWriter.Close()
	pipeWriter.Write(b)
}

func PipeReadFile(pipeReader *io.PipeReader, c chan bool) {
	if f, err := DecodeFile(pipeReader); err == nil && f.Name == testFileName {
		c <- true
	} else {
		log.Printf("Error expected nil, got %s", err)
		log.Printf("File, expected %s got %s", testFileName, f.Name)
	}
	c <- false
}

func newTestFile(fileName string) File {
	f := File{}
	f.Name = fileName
	f.Accessed = time.Now()
	f.Created = time.Now()
	f.Modified = time.Now()
	return f
}
