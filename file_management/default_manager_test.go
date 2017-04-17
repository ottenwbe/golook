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
package file_management

import (
	. "github.com/ottenwbe/golook/routing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("DefaultManager", func() {

	const (
		TEST_FOLDER = "../file_management"
		TEST_FILE   = "default_manager_test.go"
	)

	It("should emit a rpc call to the uplink server when reporting a folder and not return with an error", func() {
		RunWithMockedGolookClient(func() {
			d := defaultFileManager{}
			err := d.ReportFolder(TEST_FOLDER, false)
			Expect(err).To(BeNil())
			Expect(AccessMockedGolookClient().VisitDoPostFiles).To(BeTrue())
		})
	})

	It("should emit a rpc call to the uplink server when replacing reported files with a folder", func() {
		RunWithMockedGolookClient(func() {
			d := defaultFileManager{}
			err := d.ReportFolderR(TEST_FOLDER, false)
			Expect(err).To(BeNil())
			Expect(AccessMockedGolookClient().VisitDoPutFiles).To(BeTrue())
		})
	})

	It("should emit a rpc call to the uplink server when reporting a file", func() {
		RunWithMockedGolookClientF(func() {
			d := defaultFileManager{}
			err := d.ReportFile(TEST_FILE, false)
			Expect(err).To(BeNil())
			Expect(AccessMockedGolookClient().VisitDoPostFile).To(BeTrue())
		}, TEST_FILE, TEST_FOLDER)
	})

	It("should emit a rpc call to the uplink server when replacing a file", func() {
		RunWithMockedGolookClientF(func() {
			d := defaultFileManager{}
			err := d.ReportFileR(TEST_FILE, false)
			Expect(err).To(BeNil())
			Expect(AccessMockedGolookClient().VisitDoPutFiles).To(BeTrue())
		}, TEST_FILE, TEST_FOLDER)
	})

	It("should return an error when a folder does not exist", func() {
		RunWithMockedGolookClientF(func() {
			d := defaultFileManager{}
			err := d.ReportFolder("no_folder", false)
			Expect(err).ToNot(BeNil())
		}, TEST_FILE, TEST_FOLDER)
	})

	It("should return an error when a file does not exist which should be reported", func() {
		RunWithMockedGolookClientF(func() {
			d := defaultFileManager{}
			err := d.ReportFile("no_file.txt", false)
			Expect(err).ToNot(BeNil())
		}, TEST_FILE, TEST_FOLDER)
	})

	It("should return an error when a file does not exist which should replace the reported files", func() {
		RunWithMockedGolookClientF(func() {
			d := defaultFileManager{}
			err := d.ReportFileR("no_file.txt", false)
			Expect(err).ToNot(BeNil())
		}, TEST_FILE, TEST_FOLDER)
	})
})
