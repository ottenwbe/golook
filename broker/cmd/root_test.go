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
package cmd

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/ottenwbe/golook/broker/runtime"
	"github.com/spf13/cobra"
	"time"
)

var _ = Describe("The root command", func() {

	var (
		mockCmd = &cobra.Command{
			Use:   "mock",
			Short: "mock",
			Long:  "mock",
			Run: func(_ *cobra.Command, _ []string) {
			},
		}
	)

	It("should panic when a a call to run cannot be executed, i.e., due to parameters of a test", func() {
		tmp := RootCmd
		RootCmd = mockCmd
		Expect(func() { Run() }).To(Panic())
		RootCmd = tmp
	})

	It("should start the servers when executed.", func() {

		go RootCmd.Run(nil, nil)
		//Wait for root command in go
		time.Sleep(time.Millisecond * 600)

		Expect(len(runtime.ServerState())).To(BeNumerically(">=", 2))
		Expect(runtime.ServerState()).ToNot(ContainElement(false))

		runtime.StopServer()
	})
})
