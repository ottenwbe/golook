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
	"reflect"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/ottenwbe/golook/models"
)

var _ = Describe("The file manager's router", func() {

	Context("initialization", func() {

		var(
			testFM FileManager = nil
		)

		It("should work with a DefaultFileManager if not specified otherwise", func() {
			testFM = newFileManager()
			Expect(reflect.TypeOf(testFM)).To(Equal(reflect.TypeOf(&defaultFileManager{})))
		})
	})

	Context("routing", func() {

		const (
			TEST_PATH = "testPath"
		)

		BeforeEach(func() {
			fileManager = newTestFileManager()

		})

		AfterEach(func() {
			// reset file manager
			fileManager = newFileManager()

		})

		It("should report a file for replacement and monitoring when Monitor and Replace are true in a Filereport", func() {
			fileReport := &models.FileReport{
				Path:    TEST_PATH,
				Monitor: true,
				Replace: true,
			}

			HandleFileReport(fileReport)

			Expect(fileManager.(*testFileManager).recordedMonitor).To(BeTrue())
			Expect(fileManager.(*testFileManager).recordedReplaceMode).To(BeTrue())
			Expect(fileManager.(*testFileManager).recordedFilePath).To(Equal(TEST_PATH))
		})

		It("should report a file not for replacement but for monitoring when Monitor is true and Replace is false in a Filereport", func() {
			fileReport := &models.FileReport{
				Path:    TEST_PATH,
				Monitor: true,
				Replace: false,
			}

			HandleFileReport(fileReport)

			Expect(fileManager.(*testFileManager).recordedMonitor).To(BeTrue())
			Expect(fileManager.(*testFileManager).recordedReplaceMode).To(BeFalse())
			Expect(fileManager.(*testFileManager).recordedFilePath).To(Equal(TEST_PATH))
		})

		It("should not report a file if Filereport is nil", func() {
			var fileReport *models.FileReport = nil

			HandleFileReport(fileReport)

			Expect(fileManager.(*testFileManager).recordedMonitor).To(BeFalse())
			Expect(fileManager.(*testFileManager).recordedReplaceMode).To(BeFalse())
			Expect(fileManager.(*testFileManager).recordedFilePath).To(Equal(""))
		})

		It("should report a folder for replacement and monitoring when Monitor and Replace are true in a Filereport", func() {
			fileReport := &models.FileReport{
				Path:    TEST_PATH,
				Monitor: true,
				Replace: true,
			}

			HandleFolderReport(fileReport)

			Expect(fileManager.(*testFileManager).recordedMonitor).To(BeTrue())
			Expect(fileManager.(*testFileManager).recordedReplaceMode).To(BeTrue())
			Expect(fileManager.(*testFileManager).recordedFilePath).To(Equal(TEST_PATH))
		})

		It("should report a folder not for replacement but for monitoring when Monitor is true and Replace is false in a Filereport", func() {
			fileReport := &models.FileReport{
				Path:    TEST_PATH,
				Monitor: true,
				Replace: false,
			}

			HandleFolderReport(fileReport)

			Expect(fileManager.(*testFileManager).recordedMonitor).To(BeTrue())
			Expect(fileManager.(*testFileManager).recordedReplaceMode).To(BeFalse())
			Expect(fileManager.(*testFileManager).recordedFilePath).To(Equal(TEST_PATH))
		})

		It("should not report a folder if Filereport is nil", func() {
			var fileReport *models.FileReport = nil

			HandleFolderReport(fileReport)

			Expect(fileManager.(*testFileManager).recordedMonitor).To(BeFalse())
			Expect(fileManager.(*testFileManager).recordedReplaceMode).To(BeFalse())
			Expect(fileManager.(*testFileManager).recordedFilePath).To(Equal(""))
		})

	})
})

func newTestFileManager() FileManager {
	return &testFileManager{}
}

//File Manager which ignores all calls and records the parameters
type testFileManager struct {
	recordedFilePath    string
	recordedMonitor     bool
	recordedReplaceMode bool
}

func (fm *testFileManager) ReportFile(filePath string, monitor bool) error {
	fm.recordedFilePath = filePath
	fm.recordedMonitor = monitor
	return nil
}

func (fm *testFileManager) ReportFileR(filePath string, monitor bool) error {
	fm.recordedFilePath = filePath
	fm.recordedMonitor = monitor
	fm.recordedReplaceMode = true
	return nil
}

func (fm *testFileManager) ReportFolder(folderPath string, monitor bool) error {
	fm.recordedFilePath = folderPath
	fm.recordedMonitor = monitor
	return nil
}

func (fm *testFileManager) ReportFolderR(folderPath string, monitor bool) error {
	fm.recordedFilePath = folderPath
	fm.recordedMonitor = monitor
	fm.recordedReplaceMode = true
	return nil
}
