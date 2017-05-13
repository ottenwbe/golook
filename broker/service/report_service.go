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
	"github.com/ottenwbe/golook/broker/models"
	repo "github.com/ottenwbe/golook/broker/repository"
	"github.com/ottenwbe/golook/broker/routing"
	golook "github.com/ottenwbe/golook/broker/runtime/core"
	"github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"
)

const (
	mockReport  = "mock"
	localReport = "local"
	bCastReport = "broadcast"
)

const (
	fileReport = "file_report"
)

type (
	//reportService
	reportService interface {
		report(fileReport *models.FileReport) (map[string]map[string]*models.File, error)
		close()
	}
	//monitoredReportService is the base for all report services which monitor file changes
	monitoredReportService struct {
		fileMonitor *FileMonitor
	}
	//broadcastReportService broadcasts files to all peers
	broadcastReportService struct {
		monitoredReportService
		router           *router
		systemCallbackID string
	}
	//localReportService broadcasts files to all peers
	localReportService struct {
		monitoredReportService
	}
	/*MockReportService implements a mock version of the report service, e.g., for testing.*/
	MockReportService struct {
		FileReport *models.FileReport
	}
)

func newReportService(reportType string, router *router) (result reportService) {
	switch reportType {
	case mockReport:
		result = &MockReportService{}
	case localReport:
		result = newLocalReportService(router)
	default:
		result = newBroadcastReportService(router)
	}

	return result
}

func newBroadcastReportService(router *router) reportService {
	rs := &broadcastReportService{
		router: router,
	}

	rs.fileMonitor = &FileMonitor{}
	rs.fileMonitor.Reporter =
		func(filePath string) {
			reportFileChanges(filePath, rs.router)
		}
	rs.fileMonitor.Open()

	rs.systemCallbackID = uuid.NewV4().String()
	newSystemCallbacks.Add(
		rs.systemCallbackID,
		func(_ string, _ map[string]*golook.System) {
			broadcastLocalFiles(rs.router)
		},
	)

	router.AddHandler(fileReport,
		routing.NewHandler(
			handleFileReport,
			nil,
		),
	)

	return rs
}

func (rs *broadcastReportService) close() {
	if rs.fileMonitor != nil {
		rs.fileMonitor.Close()
	}
	newSystemCallbacks.Delete(rs.systemCallbackID)
}

func (rs *broadcastReportService) report(fileReport *models.FileReport) (map[string]map[string]*models.File, error) {

	if fileReport == nil {
		logFileReport().Error("Ignoring empty file report.")
		return map[string]map[string]*models.File{}, errors.New("Ignoring empty file report")
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

func newLocalReportService(router *router) reportService {
	rs := &localReportService{}

	rs.fileMonitor = &FileMonitor{}
	rs.fileMonitor.Reporter = reportFileChangesLocal
	rs.fileMonitor.Open()

	router.AddHandler(fileReport,
		routing.NewHandler(
			handleFileReport,
			nil,
		),
	)

	return rs
}

func (rs *localReportService) close() {
	rs.fileMonitor.Close()
}

func (rs *localReportService) report(fileReport *models.FileReport) (map[string]map[string]*models.File, error) {

	if fileReport == nil {
		logFileReport().Error("Ignoring empty file report.")
		return map[string]map[string]*models.File{}, errors.New("Ignoring empty file report")
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

func (mock *MockReportService) report(fileReport *models.FileReport) (map[string]map[string]*models.File, error) {
	mock.FileReport = fileReport
	return map[string]map[string]*models.File{}, nil
}

func (mock *MockReportService) close() {
}

func handleFileReport(params models.EncapsulatedValues) interface{} {
	var (
		fileMessage peerFileReport
		response    fileQueryData
	)

	logFileReport().WithField("params", params).Debug("New file report.")

	if err := params.Unmarshal(&fileMessage); err == nil {
		response = processFileReport(&fileMessage)
	} else {
		logFileReport().WithError(err).Error("Cannot handle file report.")
		response = fileQueryData{}
	}

	return response
}

func processFileReport(fileMessage *peerFileReport) (response fileQueryData) {
	logFileReport().Debug("Update file for: %s", fileMessage.System)
	repo.GoLookRepository.UpdateFiles(fileMessage.System, fileMessage.Files)
	return fileQueryData{}
}

func logFileReport() *log.Entry {
	return log.WithField("service", "report")
}
