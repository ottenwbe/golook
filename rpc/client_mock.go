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
package rpc

import (
	. "github.com/ottenwbe/golook/app"
	. "github.com/ottenwbe/golook/file_management"

	log "github.com/sirupsen/logrus"
)

type MockGolookClient struct {
	VisitDoPostFile  bool
	VisitDoPutFiles  bool
	VisitDoGetFiles  bool
	VisitDoPostFiles bool
	FileName         string
	FolderName       string
}

func (mock *MockGolookClient) DoPostFiles(file []File) string {
	mock.VisitDoPostFiles = true
	return ""
}

func (*MockGolookClient) DoQuerySystemsAndFiles(fileName string) (systems map[string]*System, err error) {
	panic("implement me")
}

func (*MockGolookClient) DoGetSystem(system string) (*System, error) {
	panic("implement me")
}

func (*MockGolookClient) DoPutSystem(system *System) *System {
	panic("implement me")
}

func (*MockGolookClient) DoDeleteSystem() string {
	panic("implement me")
}

func (*MockGolookClient) DoGetHome() string {
	panic("not needed")
	return ""
}

func (mock *MockGolookClient) DoPostFile(file *File) string {
	log.WithField("called", mock.VisitDoPostFile).WithField("file", *file).Info("Test DoPostFile")
	mock.VisitDoPostFile = mock.VisitDoPostFile || file != nil && file.Name == mock.FileName
	return ""
}

func (mock *MockGolookClient) DoPutFiles(files []File) string {
	mock.VisitDoPutFiles = len(files) > 0
	return ""
}

func (mock *MockGolookClient) DoGetFiles(systemName string) (map[string]File, error) {
	mock.VisitDoGetFiles = true
	return map[string]File{}, nil
}
