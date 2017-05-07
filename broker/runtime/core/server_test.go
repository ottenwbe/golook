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
	"github.com/gorilla/mux"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"sync"
	"time"
)

var _ = Describe("The server", func() {

	Context("Glob", func() {
		It("can get or create servers.", func() {
			s := NewServer(":8987", ServerHttp)
			Expect(s).ToNot(BeNil())
			Expect(servers).To(ContainElement(s))
		})
	})

	Context("http server", func() {
		It("can be started and the stopped", func() {
			s := HTTPSever{}
			s.Address = ":8765"
			s.router = mux.NewRouter().StrictSlash(true)

			var wg = sync.WaitGroup{}
			wg.Add(1)
			go s.StartServer(&wg)
			time.Sleep(time.Millisecond * 300)
			Expect(s.IsRunning()).To(BeTrue())
			s.Stop()
			Expect(s.IsRunning()).To(BeFalse())
		})
	})

})
