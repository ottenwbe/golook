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
package config

import (
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"io/ioutil"
	"os"
	"testing"
)

func TestClients(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Config Suite")
}

var _ = Describe("The main configuration", func() {
	Context("default", func() {
		It("should return the current version", func() {
			Expect(versionFromCommand()).To(Equal(VERSION))
		})
	})
})

//https://play.golang.org/p/fXpK0ZhXXf
func versionFromCommand() string {

	defer func(reset *os.File) {
		os.Stdout = reset
	}(os.Stdout)
	r, w, errPipe := os.Pipe()
	if errPipe != nil {
		return fmt.Sprintf("Pipe Error: %s", errPipe)
	}
	os.Stdout = w

	cmdVersion.Run(nil, []string{})

	w.Close()
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return fmt.Sprintf("Read Error: %s", err)
	}

	return string(b)
}
