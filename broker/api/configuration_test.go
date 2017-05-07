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

package api

import (
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	golook "github.com/ottenwbe/golook/broker/runtime/core"
	"github.com/spf13/viper"
)

var _ = Describe("The api's configuration", func() {

	Context("for the endpoints", func() {

		It("should be initialized with default values.", func() {

			viper.Reset()

			InitConfiguration()

			Expect(viper.GetBool("api.info")).To(BeTrue())
			Expect(viper.GetString("api.server.address")).To(Equal(":8383"))
		})

		It("should instantiate the Server as HTTPServer.", func() {

			ApplyConfiguration()

			httpServer, ok := HTTPServer.(*golook.HTTPSever)

			Expect(httpServer).ToNot(BeNil())
			Expect(ok).To(BeTrue())
		})

		It("should register all http end points.", func() {
			viper.SetDefault("api.info", true)

			ApplyConfiguration()

			ep := HTTPServer.(*golook.HTTPSever).RegisteredEndpoints()
			Expect(extractStringFromSlice(InfoEndpoint, ep)).To(Equal(InfoEndpoint))
			Expect(extractStringFromSlice(FileEndpoint, ep)).To(Equal(FileEndpoint))
			Expect(extractStringFromSlice(QueryEndpoint, ep)).To(Equal(QueryEndpoint))
			Expect(extractStringFromSlice(HTTPApiEndpoint, ep)).To(Equal(HTTPApiEndpoint))
			Expect(extractStringFromSlice(ConfigEndpoint, ep)).To(Equal(ConfigEndpoint))
			Expect(extractStringFromSlice(LogEndpoint, ep)).To(Equal(LogEndpoint))
			Expect(extractStringFromSlice(SystemEndpoint, ep)).To(Equal(SystemEndpoint))
		})

	})
})

func extractStringFromSlice(searchString string, slice []string) string {
	for _, ep := range slice {
		if ep == searchString {
			return ep
		}
	}
	return fmt.Sprintf("cannot find '%s' in list of endpoints", searchString)
}
