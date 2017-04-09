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
package control

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/ottenwbe/golook/routing"
)

var _ = Describe("The report service", func() {

	const (
		FILE_NAME   string = "default_test.go"
		FOLDER_NAME string = "."
	)

	var (
		ctl LookController
	)

	BeforeEach(func() {
		ctl = NewController()
	})

	It("should call the golook routing with a given file", func() {
		RunWithMockedGolookClientF(func() {
			ctl.ReportFile(FILE_NAME)
			Expect(routing.GolookClient.(*MockGolookClient).visitDoPostFile).To(BeTrue())
		}, FILE_NAME, FOLDER_NAME)
	})

	It("should NOT call the golook routing with a non existing file", func() {
		RunWithMockedGolookClient(func() {
			ctl.ReportFile(FILE_NAME + ".abc")
			Expect(routing.GolookClient.(*MockGolookClient).visitDoPostFile).To(BeFalse())
		})
	})

	It("should call the golook routing for a given folder", func() {
		RunWithMockedGolookClientF(func() {
			ctl.ReportFolder(FOLDER_NAME)
			Expect(routing.GolookClient.(*MockGolookClient).visitDoPostFiles).To(BeTrue())
		}, FILE_NAME, FOLDER_NAME)
	})

	It("should NOT call the golook routing with a non existing file", func() {
		RunWithMockedGolookClient(func() {
			ctl.ReportFolder("no_folder")
			Expect(routing.GolookClient.(*MockGolookClient).visitDoPostFile).To(BeFalse())
		})
	})

	It("should call the golook routing with files from existing folder which replace reported files", func() {
		RunWithMockedGolookClient(func() {
			ctl.ReportFolderR(FOLDER_NAME)
			Expect(routing.GolookClient.(*MockGolookClient).visitDoPutFiles).To(BeTrue())
		})
	})

	It("should NOT call the golook routing with files from existing folder when folder does not exist", func() {
		RunWithMockedGolookClient(func() {
			ctl.ReportFolderR("no_folder")
			Expect(routing.GolookClient.(*MockGolookClient).visitDoPutFiles).To(BeFalse())
		})
	})
})

var _ = Describe("The query service", func() {

	var (
		ctl LookController
	)

	BeforeEach(func() {
		ctl = NewController()
	})

	It("should call the golook routing", func() {
		RunWithMockedGolookClient(func() {
			_, err := ctl.QueryReportedFiles()
			Expect(err).To(BeNil())
			Expect(routing.GolookClient.(*MockGolookClient).visitDoGetFiles).To(BeTrue())
		})
	})
})
