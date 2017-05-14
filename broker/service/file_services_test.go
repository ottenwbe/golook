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

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/ottenwbe/golook/broker/models"
	"reflect"
)

var _ = Describe("The file services", func() {
	It("can create a file service where all files are broadcast.", func() {
		fileServices := OpenFileServices(BroadcastFiles)
		defer CloseFileServices(fileServices)
		Expect(fileServices).ToNot(BeNil())
		Expect(reflect.TypeOf(fileServices).String()).ToNot(Equal(reflect.TypeOf(scenarioBroadcastFiles{}).String()))
	})

	It("can create a file service where all queries are boradcast.", func() {
		fileServices := OpenFileServices(BroadcastQueries)
		defer CloseFileServices(fileServices)
		Expect(fileServices).ToNot(BeNil())
		Expect(reflect.TypeOf(fileServices).String()).ToNot(Equal(reflect.TypeOf(scenarioBroadcastQueries{}).String()))
	})

	It("calls the query service for queries.", func() {
		const expectedQuery = "query some file"

		fileServices := OpenFileServices(MockFileServices)
		defer CloseFileServices(fileServices)

		fileServices.Query(expectedQuery)

		Expect(AccessMockedQueryService(fileServices).SearchString).To(Equal(expectedQuery))
	})

	It("calls the report service for reports.", func() {
		var expectedFileReport = &models.FileReport{".", false}

		fileServices := OpenFileServices(MockFileServices)
		defer CloseFileServices(fileServices)

		fileServices.Report(expectedFileReport)

		Expect(AccessMockedReportService(fileServices).FileReport).To(Equal(expectedFileReport))
	})
})
