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
package data_manipulation

import (
	. "github.com/onsi/ginkgo"
	"github.com/ottenwbe/golook/client"
	"github.com/ottenwbe/golook/utils"
)

var _ = Describe("The report service", func() {
	It("should call the golook client with a given file", func() {

		client.GolookClient = testGolookClient{
			visitDoPostFile: false,
		}

		//ReportFile("reporting_test.go")

		//Expect(client.GolookClient.(testGolookClient).visitDoPostFile).To(BeTrue())
	})

})

type testGolookClient struct {
	visitDoPostFile bool
}

func (testGolookClient) DoGetHome() string {
	panic("not needed")
	return ""
}

func (t testGolookClient) DoPostFile(file *utils.File) string {
	t.visitDoPostFile = true
	return ""
}
