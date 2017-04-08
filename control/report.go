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
package control

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
func ReportFolderR(folderPath string) error {
	report, err := generateReport(folderPath)
	GolookClient.DoPutFiles(report)
	return err
}

// Report files in a folder
func ReportFolder(folderPath string) error {
	report, err := generateReport(folderPath)
	GolookClient.DoPostFiles(report)
	return err
}

func generateReport(folderPath string) ([]utils.File, error) {

	var (
		files     []os.FileInfo
		returnErr error        = nil
		report    []utils.File = make([]utils.File, 0)
	)

	files, returnErr = ioutil.ReadDir(folderPath)
	if returnErr != nil {
		return report, returnErr
	}

	for idx := range files {
		report, returnErr = appendFile(files[idx], report)
	}
	return report, returnErr
}

func appendFile(fileToAppend os.FileInfo, appendReport []utils.File) (report []utils.File, err error) {
	var file *utils.File = nil
	if file, err = utils.NewFile(fileToAppend.Name()); err == nil && !fileToAppend.IsDir() {
		report = append(appendReport, *file)
	}
	return
}
