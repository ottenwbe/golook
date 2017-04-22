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
	log "github.com/sirupsen/logrus"

	. "github.com/ottenwbe/golook/broker/models"
	. "github.com/ottenwbe/golook/broker/routing"
	. "github.com/ottenwbe/golook/broker/runtime"
	"io/ioutil"
	"os"
)

//TODO refactor

func routeFile(filePath string, replace bool) {

	file, err := NewFile(filePath)
	if err != nil {
		log.WithError(err).Error("Ignoring file api.")
		return
	}

	message, err := Marshal(FileTransfer{Files: map[string]*File{file.Name: file}, Replace: replace, System: GolookSystem.UUID})
	if err != nil {
		log.WithError(err).Error("Ignoring file api.")
		return
	}

	/* api file */
	systemIndex.Route(SysKey(), FILE_REPORT, message)
}

func routeFolder(folderPath string, replace bool) {
	report, err := generateReport(folderPath)
	if err != nil {
		log.WithError(err).Error("Ignoring folder api.")
		return
	}

	message, err := Marshal(FileTransfer{Files: report, Replace: replace, System: GolookSystem.UUID})
	if err != nil {
		log.WithError(err).Error("Ignoring folder api.")
		return
	}

	systemIndex.Route(SysKey(), FILE_REPORT, message)
}

// Generate a []File array from files in a folder
func generateReport(folderPath string) (map[string]*File, error) {

	var (
		files     []os.FileInfo
		report    map[string]*File = make(map[string]*File)
		returnErr error            = nil
	)

	files, returnErr = ioutil.ReadDir(folderPath)
	if returnErr != nil {
		return report, returnErr
	}

	for idx := range files {
		report = appendFile(files[idx], report)
	}
	return report, returnErr
}

func appendFile(fileToAppend os.FileInfo, report map[string]*File) map[string]*File {
	if file, err := NewFile(fileToAppend.Name()); err == nil && !fileToAppend.IsDir() {
		report[file.Name] = file
	}
	return report
}
