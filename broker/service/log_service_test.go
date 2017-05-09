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

package service

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	log "github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
)

var _ = Describe("The log services", func() {
	It("can read a log file from disk and write it to an io.Writer.", func() {
		LogService := GetLogService()
		log.Info("Test log entry")

		r, w := io.Pipe()
		go func() {
			LogService.RewriteLog(w)
			w.Close()
		}()

		b, err := ioutil.ReadAll(r)
		if err != nil {
			log.Fatal("Cannot read from log file in test.")
		}

		Expect(string(b)).ToNot(Equal("Test log entry"))
	})
})
