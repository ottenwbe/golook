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
	"io/ioutil"
	"os"

	. "github.com/ottenwbe/golook/routing"
	"github.com/ottenwbe/golook/utils"
)

// Report individual files
func ReportFile(filePath string) error {
	if f, err := utils.NewFile(filePath); err != nil {
		return err
	} else /* report file */ {
		GolookClient.DoPostFile(f)
	}
	return nil
}

// Report files in a folder and replace all previously reported files
func ReportFolderR(folderPath string) (err error) {
	report := make([]utils.File, 0)
	var files []os.FileInfo
	err = nil

	files, err = ioutil.ReadDir(folderPath)
	if err != nil {
		return
	}

	for idx := range files {
		err, report = reportFilesOfFolder(files[idx], report)
	}

	GolookClient.DoPutFiles(report)
	return
}

// Report files in a folder
func ReportFolder(folderPath string) (err error) {
	report := make([]utils.File, 0)
	var files []os.FileInfo
	err = nil

	files, err = ioutil.ReadDir(folderPath)
	if err != nil {
		return
	}

	for idx := range files {
		err, report = reportFilesOfFolder(files[idx], report)
	}
	GolookClient.DoPostFiles(report)
	return
}

func reportFilesOfFolder(fi os.FileInfo, oldReport []utils.File) (err error, report []utils.File) {
	var file *utils.File = nil
	if file, err = utils.NewFile(fi.Name()); err == nil && !fi.IsDir() {
		report = append(oldReport, *file)
	}
	return
}