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
	. "github.com/ottenwbe/golook/rpc"
)

var _ = Describe("The report service", func() {

	const (
		FILE_NAME   string = "default_router_test.go"
		FOLDER_NAME string = "."
	)

	var (
		ctl LookRouter
	)

	BeforeEach(func() {
		ctl = NewRouter()
	})

	It("should call the golook rpc with a given file", func() {
		RunWithMockedGolookClientF(func() {
			ctl.ReportFile(FILE_NAME)
			Expect(AccessMockedGolookClient().VisitDoPostFile).To(BeTrue())
		}, FILE_NAME, FOLDER_NAME)
	})

	It("should NOT call the golook rpc with a non existing file", func() {
		RunWithMockedGolookClient(func() {
			ctl.ReportFile(FILE_NAME + ".abc")
			Expect(AccessMockedGolookClient().VisitDoPostFile).To(BeFalse())
		})
	})

	It("should call the golook rpc for a given folder", func() {
		RunWithMockedGolookClientF(func() {
			ctl.ReportFolder(FOLDER_NAME)
			Expect(AccessMockedGolookClient().VisitDoPostFiles).To(BeTrue())
		}, FILE_NAME, FOLDER_NAME)
	})

	It("should NOT call the golook rpc with a non existing file", func() {
		RunWithMockedGolookClient(func() {
			ctl.ReportFolder("no_folder")
			Expect(AccessMockedGolookClient().VisitDoPostFile).To(BeFalse())
		})
	})

	It("should call the golook rpc with files from existing folder which replace reported files", func() {
		RunWithMockedGolookClient(func() {
			ctl.ReportFolderR(FOLDER_NAME)
			Expect(AccessMockedGolookClient().VisitDoPutFiles).To(BeTrue())
		})
	})

	It("should NOT call the golook rpc with files from existing folder when folder does not exist", func() {
		RunWithMockedGolookClient(func() {
			ctl.ReportFolderR("no_folder")
			Expect(AccessMockedGolookClient().VisitDoPutFiles).To(BeFalse())
		})
	})
})

var _ = Describe("The query service", func() {

	var (
		ctl LookRouter
	)

	BeforeEach(func() {
		ctl = NewRouter()
	})

	It("should call the golook rpc", func() {
		RunWithMockedGolookClient(func() {
			_, err := ctl.QueryReportedFiles()
			Expect(err).To(BeNil())
			Expect(AccessMockedGolookClient().VisitDoGetFiles).To(BeTrue())
		})
	})
})
