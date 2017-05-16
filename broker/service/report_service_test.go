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
	"github.com/ottenwbe/golook/broker/repository"
	. "github.com/ottenwbe/golook/broker/repository"
	"github.com/ottenwbe/golook/broker/routing"
	golook "github.com/ottenwbe/golook/broker/runtime/core"
	"github.com/ottenwbe/golook/utils"
	"path/filepath"
)

var _ = Describe("The broadcast report service", func() {

	var rs *broadcastReportService

	BeforeEach(func() {
		rs = newReportService(bCastReport, &router{routing.NewMockedRouter()}).(*broadcastReportService)
	})

	AfterEach(func() {
		rs.close()
	})

	It("ignores nil file reports", func() {
		beforeCount := routing.AccessMockedRouter(rs.router.Router.(*routing.MockRouter)).Visited
		rs.report(nil)
		Expect(routing.AccessMockedRouter(rs.router.Router.(*routing.MockRouter)).Visited).To(Equal(beforeCount))
	})

	It("broadcasts file reports and adds them to the local repository", func() {

		testFileName := "report_service_test.go"
		rs.report(
			&models.FileReport{
				Path: testFileName,
			},
		)

		storedFiles := repositories.GoLookRepository.GetFiles(golook.GolookSystem.UUID)

		Expect(len(storedFiles) >= 1).To(BeTrue())
		Expect(rs.router.Router.(*routing.MockRouter).Visited >= 1).To(BeTrue())

	})

	It("broadcasts files in folders reports and adds them to the local repository", func() {

		testFolderName, err := filepath.Abs(".")
		if err != nil {
			Fail("Cannot prepare test data, testfolder '.' not found.")
		}
		rs.report(
			&models.FileReport{
				Path: testFolderName,
			},
		)

		storedFiles := repositories.GoLookRepository.GetFiles(golook.GolookSystem.UUID)

		Expect(len(storedFiles)).To(BeNumerically(">=", 1))
		Expect(rs.router.Router.(*routing.MockRouter).Visited >= 1).To(BeTrue())

	})
})

var _ = Describe("The report handler", func() {

	BeforeEach(func() {
		//reset repo ... for each test
		GoLookRepository = NewRepository()
	})

	AfterEach(func() {
		//reset repo ... would actually be sufficient after the last test
		GoLookRepository = NewRepository()
	})

	It("stores file reports in the golook repository (when the system is known to the repo)", func() {
		storedSys := golook.GolookSystem
		GoLookRepository.StoreSystem(golook.GolookSystem.UUID, storedSys)
		f, _ := models.NewFile("report_handler_test.go")
		m, _ := utils.MarshalS(peerFileReport{Files: map[string]map[string]*models.File{filepath.Dir(f.Name): {f.Name: f}}, SystemUUID: golook.GolookSystem.UUID})
		fReport := routing.Params(m)
		handleFileReport(fReport)
		_, ok := AccessMapRepository().GetSystem(golook.GolookSystem.UUID)
		Expect(ok).To(BeTrue())
	})

	It("rejects invalid reports)", func() {
		//GoLookRepository.StoreSystem(runtime.GolookSystem.SystemUUID, runtime.GolookSystem)
		fReport := routing.Params("")
		handleFileReport(fReport)
		_, ok := AccessMapRepository().GetSystem(golook.GolookSystem.UUID)
		Expect(ok).To(BeFalse())
	})
})
