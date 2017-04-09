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
	"log"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"runtime"
)

const (
	testFileName = "test.txt"
	testSysName  = "test"
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

	Context("System", func() {
		It("can be marshalled and unmarshalled", func() {
			go PipeWriteSystem(pipeWriter)
			go PipeReadSystem(pipeReader, chanbools)
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

	Context("Systems", func() {
		It("can be marshalled and unmarshalled", func() {
			go PipeWriteSystems(pipeWriter)
			go PipeReadSystems(pipeReader, chanbools)
			result := <-chanbools
			Expect(result).To(BeTrue())
		})

		It("a new system is not nil", func() {
			Expect(NewSystem()).ToNot(BeNil())
		})

		It("a new system has no uuid be assigned", func() {
			Expect(systemUuid()).To(Equal(""))
		})

		It("can be assigned the correct os", func() {
			Expect(systemOS()).To(Equal(runtime.GOOS))
		})

	})

})

func newTestSystem(sysName string) *System {
	s := &System{
		Name:  sysName,
		OS:    "Linux",
		IP:    "localhost",
		UUID:  "uuid",
		Files: nil,
	}
	return s
}

func newTestFile(fileName string) File {
	f := File{}
	f.Name = fileName
	f.Accessed = time.Now()
	f.Created = time.Now()
	f.Modified = time.Now()
	return f
}

func PipeWriteSystem(pipeWriter *io.PipeWriter) {
	b, _ := json.Marshal(newTestSystem(testSysName))
	defer pipeWriter.Close()
	pipeWriter.Write(b)
}

func PipeReadSystem(pipeReader *io.PipeReader, c chan bool) {
	if s, err := DecodeSystem(pipeReader); err == nil && s.Name == testSysName {
		c <- true
	} else {
		log.Printf("Error expected nil, got %s", err)
		log.Printf("File, expected %s got %s", testSysName, s.Name)
	}
	c <- false
}

func PipeWriteSystems(pipeWriter *io.PipeWriter) {
	b, _ := json.Marshal([]*System{
		newTestSystem(testSysName),
	})
	defer pipeWriter.Close()
	pipeWriter.Write(b)
}

func PipeReadSystems(pipeReader *io.PipeReader, c chan bool) {
	if s, err := DecodeSystems(pipeReader); err == nil && s[0].Name == testSysName {
		c <- true
	} else {
		log.Printf("Error expected nil, got %s", err)
		log.Printf("File, expected %s got %s", testSysName, s[0].Name)
	}
	c <- false
}

func PipeWriteFiles(pipeWriter *io.PipeWriter) {
	b, _ := json.Marshal([]File{newTestFile(testFileName), newTestFile(testFileName)})
	defer pipeWriter.Close()
	pipeWriter.Write(b)
}

func PipeReadFiles(pipeReader *io.PipeReader, c chan bool) {
	if f, err := DecodeFiles(pipeReader); err == nil && f[0].Name == testFileName {
		c <- true
	} else {
		log.Printf("Error expected nil, got %s", err)
		log.Printf("File, expected %s got %s", testFileName, f[0].Name)
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
