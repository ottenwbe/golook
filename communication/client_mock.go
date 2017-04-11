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
package communication

import (
	"github.com/ottenwbe/golook/utils"

	log "github.com/sirupsen/logrus"
)

func RunWithMockedGolookClient(mockedFunction func()) {
	RunWithMockedGolookClientF(mockedFunction, "", "")
}

func RunWithMockedGolookClientF(mockedFunction func(), fileName string, folderName string) {
	//ensure that the GolookClient is reset after the function's execution
	defer func(reset LookClient) {
		GolookClient = reset
	}(GolookClient)

	//create a mock communication
	GolookClient = &MockGolookClient{
		VisitDoPostFile:  false,
		VisitDoPutFiles:  false,
		VisitDoGetFiles:  false,
		VisitDoPostFiles: false,
		fileName:         fileName,
		folderName:       folderName,
	}

	mockedFunction()
}

type MockGolookClient struct {
	VisitDoPostFile  bool
	VisitDoPutFiles  bool
	VisitDoGetFiles  bool
	VisitDoPostFiles bool
	fileName         string
	folderName       string
}

func AccessMockedGolookClient() *MockGolookClient {
	return GolookClient.(*MockGolookClient)
}

func (mock *MockGolookClient) DoPostFiles(file []utils.File) string {
	mock.VisitDoPostFiles = true
	return ""
}

func (*MockGolookClient) DoQuerySystemsAndFiles(fileName string) (systems map[string]*utils.System, err error) {
	panic("implement me")
}

func (*MockGolookClient) DoGetSystem(system string) (*utils.System, error) {
	panic("implement me")
}

func (*MockGolookClient) DoPutSystem(system *utils.System) *utils.System {
	panic("implement me")
}

func (*MockGolookClient) DoDeleteSystem() string {
	panic("implement me")
}

func (*MockGolookClient) DoGetHome() string {
	panic("not needed")
	return ""
}

func (mock *MockGolookClient) DoPostFile(file *utils.File) string {
	log.WithField("called", mock.VisitDoPostFile).WithField("file", *file).Info("Test DoPostFile")
	mock.VisitDoPostFile = mock.VisitDoPostFile || file != nil && file.Name == mock.fileName
	return ""
}

func (mock *MockGolookClient) DoPutFiles(files []utils.File) string {
	mock.VisitDoPutFiles = len(files) > 0
	return ""
}

func (mock *MockGolookClient) DoGetFiles() ([]utils.File, error) {
	mock.VisitDoGetFiles = true
	return []utils.File{}, nil
}
