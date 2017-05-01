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

type FileServiceType string

const (
	MockFileServices FileServiceType = "mock"
	BroadcastFiles   FileServiceType = "fileBoroadcast"
	BroadcastQueries FileServiceType = "queryBroadcast"
)

var (
	ReportService reportService
	QueryService  queryService
)

func OpenFileServices(fileServiceType FileServiceType) {
	ReportService, QueryService = newFileServices(fileServiceType)
}

func CloseFileServices() {
	if ReportService != nil {
		ReportService.Close()
	}
	ReportService = nil
	QueryService = nil
}

func newFileServices(fileServiceType FileServiceType) (reportService, queryService) {
	switch fileServiceType {
	case MockFileServices:
		return newReportService(MockReport), newQueryService(MockQueries)
	case BroadcastQueries:
		return newReportService(LocalReport), newQueryService(BCastQueries)
	default: /*BroadcastFiles*/
		return newReportService(BCastReport), newQueryService(LocalQueries)
	}
}
