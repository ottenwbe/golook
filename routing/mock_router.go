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
	"github.com/ottenwbe/golook/utils"
	log "github.com/sirupsen/logrus"
)

type MockedLookController struct {
	Visited bool
}

func NewMockedRouter() LookRouter {
	return &MockedLookController{}
}

func (mlc *MockedLookController) QueryAllSystemsForFile(fileName string) (files map[string]*utils.System, err error) {
	log.Debug("Mocked query for all Systems for file.")
	mlc.Visited = true
	return nil, nil
}

func (mlc *MockedLookController) QueryReportedFiles() (files []utils.File, err error) {
	log.Debug("Mocked reported files.")
	mlc.Visited = true
	return []utils.File{}, nil
}

func (mlc *MockedLookController) QueryFiles(systemName string) (files []utils.File, err error) {
	mlc.Visited = true
	return []utils.File{}, nil
}

func (mlc *MockedLookController) ReportFile(filePath string) error {
	mlc.Visited = true
	return nil
}

func (mlc *MockedLookController) ReportFolderR(folderPath string) error {
	mlc.Visited = true
	return nil
}

func (mlc *MockedLookController) ReportFolder(folderPath string) error {
	mlc.Visited = true
	return nil
}
