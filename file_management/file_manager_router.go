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
package file_management

import (
	. "github.com/ottenwbe/golook/models"
	log "github.com/sirupsen/logrus"
)

var (
	fileManager = newFileManager()
)

func HandleFileReport(fileReport *FileReport) {

	if fileReport == nil {
		log.Error("Ignoring empty file report ")
		return
	}

	if fileReport.Replace {
		fileManager.ReportFileR(fileReport.Path, fileReport.Monitor)
	} else {
		fileManager.ReportFile(fileReport.Path, fileReport.Monitor)
	}
}

func HandleFolderReport(folderReport *FileReport) {

	if folderReport == nil {
		log.Error("Ignoring empty file report ")
		return
	}

	if folderReport.Replace {
		fileManager.ReportFileR(folderReport.Path, folderReport.Monitor)
	} else {
		fileManager.ReportFile(folderReport.Path, folderReport.Monitor)
	}
}

func newFileManager() FileManager {
	return &defaultFileManager{}
}