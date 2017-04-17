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
package routing

import (
	. "github.com/ottenwbe/golook/app"
	. "github.com/ottenwbe/golook/models"
	. "github.com/ottenwbe/golook/repository"
)

type DefaultRouter struct {
}

func (*DefaultRouter) handleQueryAllSystemsForFile(fileName string) (systems map[string]*System, err error) {
	if canAnswerQuery() {
		systems = GoLookRepository.FindSystemAndFiles(fileName)
	} else {
		systems, err = golookClient.DoQuerySystemsAndFiles(fileName)
	}
	return
}

func (*DefaultRouter) handleQueryFiles(systemName string) (files map[string]File, err error) {
	if canAnswerQuery() {
		files, _ = GoLookRepository.GetFilesOfSystem(systemName)
	} else {
		files, err = golookClient.DoGetFiles(systemName)
	}
	return
}

func (*DefaultRouter) handleReportFile(system string, filePath string) error {
	if GolookSystem.UUID == system {
		//return DefaultFileManager{}.ReportFile(filePath, false)
	} else {
		//Route Files to uplink
		//golookClient.DoPostFile()
	}
	return nil
}

// Report individual files
func (*DefaultRouter) handleReportFileR(system string, filePath string) error {
	if GolookSystem.UUID == system {
		//return DefaultFileManager{}.ReportFileR(filePath, false)
	} else {
		//Register System with uplink
		// golookClient.Do
	}
	return nil
}

func (*DefaultRouter) handleReportFolderR(system string, folderPath string) error {
	if GolookSystem.UUID == system {
		//return DefaultFileManager{}.ReportFolderR(folderPath, false)
	} else {
		//Register System with uplink
		// golookClient.Do
	}
	return nil
}

// Report files in a folder and replace all previously reported files
func (*DefaultRouter) handleReportFolder(system string, folderPath string) error {
	if GolookSystem.UUID == system {
		//return DefaultFileManager{}.ReportFolder(folderPath, false)
	} else {
		//Register System with uplink
		// golookClient.Do
	}
	return nil
}

func canAnswerQuery() bool {
	return golookClient == nil && GoLookRepository != nil
}
