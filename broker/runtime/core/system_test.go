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

package core

import (
	//"encoding/json"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"io"
	"runtime"
	//log "github.com/sirupsen/logrus"
)

const (
	testSysName = "test"
)

var _ = Describe("The system model", func() {

	var (
		chanbools  chan bool
		pipeReader *io.PipeReader
		pipeWriter *io.PipeWriter
	)

	BeforeEach(func() {
		chanbools = make(chan bool)
		pipeReader, pipeWriter = io.Pipe()
	})

	Context("New system", func() {
		It("a new system is not nil", func() {
			Expect(NewSystem()).ToNot(BeNil())
		})

		It("has a uuid be assigned", func() {
			Expect(getUUID()).ToNot(Equal(""))
		})

		It("can be assigned the correct os", func() {
			Expect(getOS()).To(Equal(runtime.GOOS))
		})
	})

	/*Context("System conversion", func() {
		It("can be marshalled and unmarshalled", func() {
			go PipeWriteSystem(pipeWriter)
			go PipeReadSystem(pipeReader, chanbools)
			result := <-chanbools
			Expect(result).To(BeTrue())
		})

		It("an array can be marshalled and unmarshalled", func() {
			go PipeWriteSystems(pipeWriter)
			go PipeReadSystems(pipeReader, chanbools)
			result := <-chanbools
			Expect(result).To(BeTrue())
		})
	})*/
})

func newTestSystem(sysName string) *System {
	s := &System{
		Name: sysName,
		OS:   "Linux",
		IP:   "localhost",
		UUID: "uuid",
	}
	return s
}

/*
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
*/
