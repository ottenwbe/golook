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
package management

//
import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/ottenwbe/golook/broker/models"
	"github.com/ottenwbe/golook/broker/routing"
	"path/filepath"
)

var _ = Describe("The file service", func() {

	It("ignores nil file reports", func() {
		routing.RunWithMockedRouter(func() {
			MakeFileReport(nil)
			Expect(routing.AccessMockedRouter().Visited).To(BeZero())
		})
	})

	It("adds files which sepecify a monitoring flag to the file monitor", func() {
		routing.RunWithMockedRouter(func() {
			testFileName := "test_add_remove.txt"
			MakeFileReport(
				&models.FileReport{
					Path:    testFileName,
					Monitor: true,
					Replace: true,
				},
			)
			defer RemoveFileMonitor(testFileName)

			Expect(watchedFiles[testFileName]).To(BeTrue())
			Expect(routing.AccessMockedRouter().Visited >= 1).To(BeTrue())
		})
	})

	It("does not add file reports without monitoring flag to the file monitor", func() {
		routing.RunWithMockedRouter(func() {
			testFileName := "test_add_remove.txt"
			MakeFileReport(
				&models.FileReport{
					Path:    testFileName,
					Monitor: false,
					Replace: true,
				},
			)
			defer RemoveFileMonitor(testFileName)

			_, ok := watchedFiles[testFileName]
			Expect(ok).To(BeFalse())
			Expect(routing.AccessMockedRouter().Visited >= 1).To(BeTrue())
		})
	})

	It("ignores nil folder reports", func() {
		routing.RunWithMockedRouter(func() {
			MakeFolderReport(nil)
			Expect(routing.AccessMockedRouter().Visited).To(BeZero())
		})
	})

	It("does add folders specifying the monitor flag to the monitor and ignores invalid folders", func() {
		routing.RunWithMockedRouter(func() {
			folderName := "test_add_remove"
			MakeFolderReport(
				&models.FileReport{
					Path:    folderName,
					Monitor: true,
					Replace: true,
				},
			)
			defer RemoveFileMonitor(folderName)

			_, ok := watchedFiles[folderName]
			Expect(ok).To(BeTrue())
			Expect(routing.AccessMockedRouter().Visited).To(BeZero())
		})
	})

	It("does add folders specifying the monitor flag to the file monitor", func() {
		routing.RunWithMockedRouter(func() {
			folderName, _ := filepath.Abs(".")
			MakeFolderReport(
				&models.FileReport{
					Path:    folderName,
					Monitor: false,
					Replace: true,
				},
			)
			defer RemoveFileMonitor(folderName)

			_, ok := watchedFiles[folderName]
			Expect(ok).To(BeFalse())
			Expect(routing.AccessMockedRouter().Visited).To(Equal(1))
		})
	})

})
