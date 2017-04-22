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
package service

//
import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/ottenwbe/golook/broker/models"
	"github.com/ottenwbe/golook/broker/routing"
	"path/filepath"
)

var _ = Describe("The report service", func() {

	var rs = NewReportService()

	It("ignores nil file reports", func() {
		routing.RunWithMockedRouter(systemIndex, func() {
			systemIndex = routing.NewMockedRouter()

			rs.MakeFileReport(nil)
			Expect(routing.AccessMockedRouter(systemIndex).Visited).To(BeZero())
		})
	})

	It("adds files which sepecify a monitoring flag to the file monitor", func() {
		routing.RunWithMockedRouter(systemIndex, func() {
			systemIndex = routing.NewMockedRouter()

			testFileName := "test_add_remove.txt"
			rs.MakeFileReport(
				&models.FileReport{
					Path:    testFileName,
					Monitor: true,
					Replace: true,
				},
			)
			defer RemoveFileMonitor(testFileName)

			Expect(watchedFiles[testFileName]).To(BeTrue())
			Expect(routing.AccessMockedRouter(systemIndex).Visited >= 1).To(BeTrue())
		})
	})

	It("does not add file reports without monitoring flag to the file monitor", func() {
		routing.RunWithMockedRouter(systemIndex, func() {
			systemIndex = routing.NewMockedRouter()

			testFileName := "test_add_remove.txt"
			rs.MakeFileReport(
				&models.FileReport{
					Path:    testFileName,
					Monitor: false,
					Replace: true,
				},
			)
			defer RemoveFileMonitor(testFileName)

			_, ok := watchedFiles[testFileName]
			Expect(ok).To(BeFalse())
			Expect(routing.AccessMockedRouter(systemIndex).Visited >= 1).To(BeTrue())
		})
	})

	It("ignores nil folder reports", func() {
		routing.RunWithMockedRouter(systemIndex, func() {
			systemIndex = routing.NewMockedRouter()

			rs.MakeFolderReport(nil)
			Expect(routing.AccessMockedRouter(systemIndex).Visited).To(BeZero())
		})
	})

	It("does add folders specifying the monitor flag to the monitor and ignores invalid folders", func() {
		routing.RunWithMockedRouter(systemIndex, func() {
			systemIndex = routing.NewMockedRouter()

			folderName := "test_add_remove"
			rs.MakeFolderReport(
				&models.FileReport{
					Path:    folderName,
					Monitor: true,
					Replace: true,
				},
			)
			defer RemoveFileMonitor(folderName)

			_, ok := watchedFiles[folderName]
			Expect(ok).To(BeTrue())
			Expect(routing.AccessMockedRouter(systemIndex).Visited).To(BeZero())
		})
	})

	It("does add folders specifying the monitor flag to the file monitor", func() {
		routing.RunWithMockedRouter(systemIndex, func() {
			systemIndex = routing.NewMockedRouter()

			folderName, _ := filepath.Abs(".")
			rs.MakeFolderReport(
				&models.FileReport{
					Path:    folderName,
					Monitor: false,
					Replace: true,
				},
			)
			defer RemoveFileMonitor(folderName)

			_, ok := watchedFiles[folderName]
			Expect(ok).To(BeFalse())
			Expect(routing.AccessMockedRouter(systemIndex).Visited).To(Equal(1))
		})
	})

})
