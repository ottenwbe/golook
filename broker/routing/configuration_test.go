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

package routing

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/spf13/viper"
)

var _ = Describe("The configuration of the routing layer", func() {

	It("initializes ALL expected default values.", func() {
		InitConfiguration()

		Expect(viper.GetStringSlice(peers)).To(Equal([]string{""}))
		Expect(viper.GetInt(duplicateLength)).To(Equal(100))
	})

	It("applies ALL expected default values.", func() {
		InitConfiguration()
		ApplyConfiguration()

		Expect(defaultPeers).To(Equal([]string{""}))
		Expect(maxDuplicateMapLen).To(Equal(100))
	})
})
