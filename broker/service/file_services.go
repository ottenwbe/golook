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
	"github.com/ottenwbe/golook/broker/models"
	"github.com/ottenwbe/golook/broker/routing"
)

type FileServiceType string

const (
	MockFileServices FileServiceType = "mock"
	BroadcastFiles   FileServiceType = "fileBoroadcast"
	BroadcastQueries FileServiceType = "queryBroadcast"
)

var (
	FileServices fileServices
)

func OpenFileServices(fileServiceType FileServiceType) {
	FileServices = newFileServices(fileServiceType)
	FileServices.open()
}

func CloseFileServices() {
	if FileServices != nil {
		FileServices.close()
	}
}

func newFileServices(fileServiceType FileServiceType) fileServices {
	switch fileServiceType {
	case MockFileServices:
		return &scenarioMock{}
	case BroadcastQueries:
		return &scenarioBroadcastQueries{}
	default: /*BroadcastFiles*/
		return &scenarioBroadcastFiles{}
	}
}

type fileServices interface {
	open()
	close()
	Query(searchString string) (interface{}, error)
	Report(fileReport *models.FileReport) (map[string]*models.File, error)
}

type compoundFileServices struct {
	ReportService reportService
	QueryService  queryService
}

type scenarioBroadcastFiles struct {
	compoundFileServices
	r *Router
}

type scenarioBroadcastQueries struct {
	compoundFileServices
	r *Router
}
type scenarioMock struct {
	compoundFileServices
}

func (s *scenarioBroadcastQueries) open() {
	s.r = newRouter("broadcast_queries", routing.BroadcastRouter)
	s.ReportService = newReportService(LocalReport, s.r)
	s.QueryService = newQueryService(BCastQueries, s.r)

	s.r.AddHandler(
		fileQuery,
		routing.NewHandler(
			handleFileQuery,
			mergeFileQuery,
		),
	)
	s.r.AddHandler(fileReport,
		routing.NewHandler(
			handleFileReport,
			nil,
		),
	)
}

func (s *scenarioBroadcastQueries) close() {
	s.ReportService.close()
	s.r.close()
}

func (s *compoundFileServices) Query(searchString string) (interface{}, error) {
	return s.QueryService.query(searchString)
}

func (s *compoundFileServices) Report(fileReport *models.FileReport) (map[string]*models.File, error) {
	return s.ReportService.report(fileReport)
}

func (s *scenarioBroadcastFiles) open() {
	s.r = newRouter("broadcast_files", routing.BroadcastRouter)
	s.ReportService = newReportService(BCastReport, s.r)
	s.QueryService = newQueryService(LocalQueries, s.r)

	s.r.AddHandler(
		fileQuery,
		routing.NewHandler(
			handleFileQuery,
			mergeFileQuery,
		),
	)
	s.r.AddHandler(fileReport,
		routing.NewHandler(
			handleFileReport,
			nil,
		),
	)
}

func (s *scenarioBroadcastFiles) close() {
	s.ReportService.close()
	s.r.close()
}

func (s *scenarioMock) open() {
	s.ReportService = newReportService(MockReport, nil)
	s.QueryService = newQueryService(MockQueries, nil)
}

func (s *scenarioMock) close() {
	s.ReportService.close()
}

//TODO: defensive prog.

func AccessMockedQueryService(services fileServices) *MockQueryService {
	mockScenario := services.(*scenarioMock)
	mockService := mockScenario.QueryService.(*MockQueryService)
	return mockService
}

func AccessMockedReportService(services fileServices) *MockReportService {
	mockScenario := services.(*scenarioMock)
	mockService := mockScenario.ReportService.(*MockReportService)
	return mockService
}
