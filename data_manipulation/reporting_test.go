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
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/ottenwbe/golook/client"
	"github.com/ottenwbe/golook/utils"
	"github.com/sirupsen/logrus"
)

var _ = Describe("The report service", func() {
	It("should call the golook client with a given file", func() {
		runWithMockedGolookClient(func() {
			ReportFile(FILE_NAME)
			Expect(client.GolookClient.(*MockGolookClient).visitDoPostFile).To(BeTrue())
		})
	})

	It("should NOT call the golook client with a non existing file", func() {
		runWithMockedGolookClient(func() {
			ReportFile(FILE_NAME + ".abc")
			Expect(client.GolookClient.(*MockGolookClient).visitDoPostFile).To(BeFalse())
		})
	})

	It("should call the golook client with files from existing folder", func() {
		runWithMockedGolookClient(func() {
			ReportFolder(FOLDER_NAME)
			//Expect(client.GolookClient.(*MockGolookClient).visitDoPostFile).To(BeFalse())
		})
	})

})

const FILE_NAME = "reporting_test.go"
const FOLDER_NAME = "../data_manipulation"

func runWithMockedGolookClient(mockedFunction func()) {

	//ensure that the GolookClient is reset after the function's execution
	defer func(reset client.LookClientFunctions) {
		client.GolookClient = reset
	}(client.GolookClient)

	//create a mock client
	client.GolookClient = &MockGolookClient{
		visitDoPostFile: false,
	}

	mockedFunction()
}

type MockGolookClient struct {
	visitDoPostFile bool
	file            *utils.File
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

func (t *MockGolookClient) DoPostFile(file *utils.File) string {
	logrus.Info("Test posting")
	t.visitDoPostFile = file != nil && file.Name == FILE_NAME
	return ""
}

func (t *MockGolookClient) DoPutFiles(file []utils.File) string {
	t.visitDoPostFile = true
	return ""
}
