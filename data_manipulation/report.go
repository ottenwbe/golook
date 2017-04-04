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
package data_manipulation

import (
	. "github.com/ottenwbe/golook/client"
	"github.com/ottenwbe/golook/utils"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
)

// Report individual files
func ReportFile(filePath string) {
	if f, err := utils.NewFile(filePath); err != nil {
		log.WithError(err).Error("Could not report file")
	} else /* report file */ {
		GolookClient.DoPostFile(f)
	}
}

// Report files in a folder and replace all previously reported files
func ReportFolderR(folderPath string) {
	report := make([]utils.File, 0)

	//TODO error handling
	files, _ := ioutil.ReadDir(folderPath)
	for idx := range files {
		if file, err := utils.NewFile(files[idx].Name()); err != nil {
			log.WithError(err).Error("Could not report file")
		} else if !files[idx].IsDir() /* report file */ {
			report = append(report, *file)
		}
	}
	GolookClient.DoPutFiles(report)
}

// Report files in a folder
func ReportFolder(folderPath string) {
	//TODO error handling
	files, _ := ioutil.ReadDir(folderPath)
	for idx := range files {
		if file, err := utils.NewFile(files[idx].Name()); err != nil {
			log.WithError(err).Error("Could not report file")
		} else if !files[idx].IsDir() /* report file */ {
			GolookClient.DoPostFile(file)
		}
	}
}
