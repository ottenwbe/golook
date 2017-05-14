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

/*
FileServiceType describes the file service.
*/
type FileServiceType string

/*
Valid FileServiceTypes
*/
const (
	MockFileServices FileServiceType = "mock"
	BroadcastFiles   FileServiceType = "fileBoroadcast"
	BroadcastQueries FileServiceType = "queryBroadcast"
)

/*
OpenFileServices opens a file service of the given type. The default is BroadcastFiles.
*/
func OpenFileServices(fileServiceType FileServiceType) FileServices {
	fileServices := newFileServices(fileServiceType)
	fileServices.open()
	return fileServices
}

/*
CloseFileServices a given file service
*/
func CloseFileServices(fileServices FileServices) {
	if fileServices != nil {
		fileServices.close()
	}
}

func newFileServices(fileServiceType FileServiceType) FileServices {
	switch fileServiceType {
	case MockFileServices:
		return &scenarioMock{}
	case BroadcastQueries:
		return &scenarioBroadcastQueries{}
	default: /*BroadcastFiles*/
		return &scenarioBroadcastFiles{}
	}
}

/*
FileServices is the base interface for all file services
*/
type FileServices interface {
	open()
	close()
	Query(searchString string) (interface{}, error)
	Report(fileReport *models.FileReport) (map[string]map[string]*models.File, error)
}

type compoundFileServices struct {
	ReportService reportService
	QueryService  queryService
}

type scenarioBroadcastFiles struct {
	compoundFileServices
	r *router
}

type scenarioBroadcastQueries struct {
	compoundFileServices
	r *router
}
type scenarioMock struct {
	compoundFileServices
}

func (s *scenarioBroadcastQueries) open() {
	s.r = newRouter("broadcast_queries", routing.BroadcastRouterType)
	s.ReportService = newReportService(localReport, s.r)
	s.QueryService = newQueryService(bCastQueries, s.r)
}

func (s *scenarioBroadcastQueries) close() {
	s.ReportService.close()
	s.r.close()
}

/*
Query the file service for a specific file
*/
func (s *compoundFileServices) Query(searchString string) (interface{}, error) {
	return s.QueryService.query(searchString)
}

/*
Report a file/folder to the file service
*/
func (s *compoundFileServices) Report(fileReport *models.FileReport) (map[string]map[string]*models.File, error) {
	return s.ReportService.report(fileReport)
}

func (s *scenarioBroadcastFiles) open() {
	s.r = newRouter("broadcast_files", routing.BroadcastRouterType)
	s.ReportService = newReportService(bCastReport, s.r)
	s.QueryService = newQueryService(localQueries, s.r)
}

func (s *scenarioBroadcastFiles) close() {
	s.ReportService.close()
	s.r.close()
}

func (s *scenarioMock) open() {
	s.ReportService = newReportService(mockReport, nil)
	s.QueryService = newQueryService(mockQueries, nil)
}

func (s *scenarioMock) close() {
	s.ReportService.close()
}

//TODO: defensive prog.

/*
AccessMockedQueryService returns an encapsulated MockQueryService, iff services is a MockFileService.
*/
func AccessMockedQueryService(services FileServices) *MockQueryService {
	if mockScenario, ok := services.(*scenarioMock); ok {
		mockService := mockScenario.QueryService.(*MockQueryService)
		return mockService
	}
	return nil
}

/*
AccessMockedReportService returns an encapsulated MockReportService, iff services is a MockFileService.
*/
func AccessMockedReportService(services FileServices) *MockReportService {
	if mockScenario, ok := services.(*scenarioMock); ok {
		mockService := mockScenario.ReportService.(*MockReportService)
		return mockService
	}
	return nil
}
