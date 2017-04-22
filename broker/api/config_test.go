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

	"github.com/ottenwbe/golook/broker/runtime"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("The api configuration", func() {

	Context("Endpoints", func() {

		It("should register all http end points", func() {

			ConfigApi()

			ep := runtime.HttpServer.RegisteredEndpoints()
			Expect(extractStringFromSlice(INFO_EP, ep)).To(Equal(INFO_EP))
			Expect(extractStringFromSlice(FILE_EP, ep)).To(Equal(FILE_EP))
			Expect(extractStringFromSlice(FILE_QUERY_EP, ep)).To(Equal(FILE_QUERY_EP))
			Expect(extractStringFromSlice(FOLDER_EP, ep)).To(Equal(FOLDER_EP))
		})
	})
})

func extractStringFromSlice(searchString string, slice []string) string {
	for _, ep := range slice {
		if ep == searchString {
			return ep
		}
	}
	return fmt.Sprintf("array does not hold %s", searchString)
}
