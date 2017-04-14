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
)

var _ = Describe("The mocked router", func() {

	var (
		router *MockedLookRouter
	)

	BeforeEach(func() {
		router = NewMockedRouter().(*MockedLookRouter)
	})

	It("should set the valid flag to true after calling handleQueryAllSystemsForFile", func() {
		router.handleQueryAllSystemsForFile("test.txt")
		Expect(router.Visited).To(BeTrue())
	})

	It("should set the valid flag to true after calling QueryFiles", func() {
		router.handleQueryFiles("system")
		Expect(router.Visited).To(BeTrue())
	})

	It("should set the valid flag to true after calling ReportFolder", func() {
		router.handleReportFolder("folderName")
		Expect(router.Visited).To(BeTrue())
	})

	It("should set the valid flag to true after calling ReportFolderR", func() {
		router.handleReportFolderR("folderName")
		Expect(router.Visited).To(BeTrue())
	})

	It("should set the valid flag to true after calling ReportFileR", func() {
		router.handleReportFileR("file.txt")
		Expect(router.Visited).To(BeTrue())
	})

	It("should set the valid flag to true after calling ReportFile", func() {
		router.handleReportFile("file.txt")
		Expect(router.Visited).To(BeTrue())
	})
})
