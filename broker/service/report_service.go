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
	"errors"
	. "github.com/ottenwbe/golook/broker/models"
	"github.com/ottenwbe/golook/broker/runtime"
	"github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"
)

const (
	MockReport  = "mock"
	LocalReport = "local"
	BCastReport = "broadcast"
)

type (
	//reportService
	reportService interface {
		report(fileReport *FileReport) (map[string]*File, error)
		close()
	}
	//monitoredReportService is the base for all report services which monitor file changes
	monitoredReportService struct {
		fileMonitor *FileMonitor
	}
	//broadcastReportService broadcasts files to all peers
	broadcastReportService struct {
		monitoredReportService
		router           *Router
		systemCallbackId string
	}
	//localReportService broadcasts files to all peers
	localReportService struct {
		monitoredReportService
	}
	MockReportService struct {
		FileReport *FileReport
	}
)

func newReportService(reportType string, router *Router) (result reportService) {
	switch reportType {
	case MockReport:
		result = &MockReportService{}
	case LocalReport:
		result = newLocalReportService()
	default:
		result = newBroadcastReportService(router)
	}

	return result
}

func newBroadcastReportService(router *Router) reportService {
	rs := &broadcastReportService{
		router: router,
	}

	rs.fileMonitor = &FileMonitor{}
	rs.fileMonitor.reporter =
		func(filePath string) {
			reportFileChanges(filePath, rs.router)
		}
	rs.fileMonitor.Open()

	rs.systemCallbackId = uuid.NewV4().String()
	newSystemCallbacks.Add(
		rs.systemCallbackId,
		func(_ string, _ map[string]*runtime.System) {
			broadcastLocalFiles(rs.router)
		},
	)

	return rs
}

func (rs *broadcastReportService) close() {
	if rs.fileMonitor != nil {
		rs.fileMonitor.Close()
	}
	newSystemCallbacks.Delete(rs.systemCallbackId)
}

func (rs *broadcastReportService) report(fileReport *FileReport) (map[string]*File, error) {

	if fileReport == nil {
		log.Error("Ignoring empty file report.")
		return map[string]*File{}, errors.New("Ignoring empty file report")
	}

	files := localFileReport(fileReport.Path, fileReport.Delete)
	broadcastFiles(files, rs.router)

	if fileReport.Delete {
		rs.fileMonitor.RemoveMonitored(fileReport.Path)
	} else {
		rs.fileMonitor.Monitor(fileReport.Path)
	}

	return files, nil
}

func newLocalReportService() reportService {
	rs := &localReportService{}

	rs.fileMonitor = &FileMonitor{}
	rs.fileMonitor.reporter = reportFileChangesLocal
	rs.fileMonitor.Open()

	return rs
}

func (rs *localReportService) close() {
	rs.fileMonitor.Close()
}

func (rs *localReportService) report(fileReport *FileReport) (map[string]*File, error) {

	if fileReport == nil {
		log.Error("Ignoring empty file report.")
		return map[string]*File{}, errors.New("Ignoring empty file report")
	}

	// initial report
	files := localFileReport(fileReport.Path, fileReport.Delete)

	// continuous report
	if fileReport.Delete {
		rs.fileMonitor.RemoveMonitored(fileReport.Path)
	} else {
		rs.fileMonitor.Monitor(fileReport.Path)
	}

	return files, nil
}

func (mock *MockReportService) report(fileReport *FileReport) (map[string]*File, error) {
	mock.FileReport = fileReport
	return map[string]*File{}, nil
}

func (mock *MockReportService) close() {
}
