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
	"github.com/ottenwbe/golook/broker/routing"
	golook "github.com/ottenwbe/golook/broker/runtime/core"
)

var _ = Describe("The broadcast report service", func() {

	var rs *broadcastReportService

	BeforeEach(func() {
		rs = newReportService(BCastReport, &router{routing.NewMockedRouter()}).(*broadcastReportService)
	})

	AfterEach(func() {
		rs.close()
	})

	It("ignores nil file reports", func() {

		rs.report(nil)
		Expect(routing.AccessMockedRouter(rs.router.Router.(*routing.MockRouter)).Visited).To(BeZero())

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

})
