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
	"io/ioutil"
	"os"

	. "github.com/ottenwbe/golook/communication"
	. "github.com/ottenwbe/golook/utils"
)

type DefaultRouter struct {
}

func (DefaultRouter) QueryAllSystemsForFile(fileName string) (systems map[string]*System, err error) {
	systems, err = GolookClient.DoQuerySystemsAndFiles(fileName)
	return
}

func (DefaultRouter) QueryReportedFiles() (files []File, err error) {
	files, err = GolookClient.DoGetFiles()
	return
}

func (DefaultRouter) QueryFiles(systemName string) (files []File, err error) {
	files, err = GolookClient.DoGetFiles()
	return
}

// Report individual files
func (DefaultRouter) ReportFile(filePath string) error {
	if f, err := NewFile(filePath); err != nil {
		return err
	} else /* report file */ {
		GolookClient.DoPostFile(f)
	}
	return nil
}

// Report files in a folder and replace all previously reported files
func (DefaultRouter) ReportFolderR(folderPath string) error {
	report, err := generateReport(folderPath)
	GolookClient.DoPutFiles(report)
	return err
}

// Report files in a folder
func (DefaultRouter) ReportFolder(folderPath string) error {
	report, err := generateReport(folderPath)
	GolookClient.DoPostFiles(report)
	return err
}

func generateReport(folderPath string) ([]File, error) {

	var (
		files     []os.FileInfo
		report    []File = make([]File, 0)
		returnErr error  = nil
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

func appendFile(fileToAppend os.FileInfo, appendReport []File) (report []File, err error) {
	var file *File = nil
	if file, err = NewFile(fileToAppend.Name()); err == nil && !fileToAppend.IsDir() {
		report = append(appendReport, *file)
	}
	return
}
