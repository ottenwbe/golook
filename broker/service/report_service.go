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
		MakeFileReport(fileReport *FileReport) (map[string]*File, error)
		Close()
	}
	//monitoredReportService is the base for all report services which monitor file changes
	monitoredReportService struct {
		fileMonitor *FileMonitor
	}
	//broadcastReportService broadcasts files to all peers
	broadcastReportService struct {
		monitoredReportService
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

func newReportService(reportType string) (result reportService) {
	switch reportType {
	case MockReport:
		result = &MockReportService{}
	case LocalReport:
		result = newLocalReportService()
	default:
		result = newBroadcastReportService()
	}

	return result
}

func newBroadcastReportService() reportService {
	rs := &broadcastReportService{}

	rs.fileMonitor = &FileMonitor{}
	rs.fileMonitor.reporter = reportFileChanges
	rs.fileMonitor.Start()

	rs.systemCallbackId = uuid.NewV4().String()
	newSystemCallbacks.Add(
		rs.systemCallbackId,
		func(_ string, _ *runtime.System) {
			broadcastLocalFiles()
		},
	)

	return rs
}

func (rs *broadcastReportService) Close() {
	rs.fileMonitor.Close()
	newSystemCallbacks.Delete(rs.systemCallbackId)
}

func (rs *broadcastReportService) MakeFileReport(fileReport *FileReport) (map[string]*File, error) {

	if fileReport == nil {
		log.Error("Ignoring empty file report.")
		return map[string]*File{}, errors.New("Ignoring empty file report")
	}

	files := localFileReport(fileReport.Path)
	broadcastFiles(files)

	rs.fileMonitor.Monitor(fileReport.Path)

	return files, nil
}

func newLocalReportService() reportService {
	rs := &localReportService{}

	rs.fileMonitor = &FileMonitor{}
	rs.fileMonitor.reporter = reportFileChangesLocal
	rs.fileMonitor.Start()

	return rs
}

func (rs *localReportService) Close() {
	rs.fileMonitor.Close()
}

func (rs *localReportService) MakeFileReport(fileReport *FileReport) (map[string]*File, error) {

	if fileReport == nil {
		log.Error("Ignoring empty file report.")
		return map[string]*File{}, errors.New("Ignoring empty file report")
	}

	// initial report
	files := localFileReport(fileReport.Path)

	// continous report
	rs.fileMonitor.Monitor(fileReport.Path)

	return files, nil
}

func (mock *MockReportService) MakeFileReport(fileReport *FileReport) (map[string]*File, error) {
	mock.FileReport = fileReport
	return map[string]*File{}, nil
}

func (mock *MockReportService) Close() {

}
