package data_manipulation

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/ottenwbe/golook/client"
	"github.com/ottenwbe/golook/utils"
	"github.com/sirupsen/logrus"
	"testing"
)

func TestClients(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Data Manipulation Test Suite")
}

const FILE_NAME = "reporting_test.go"
const FOLDER_NAME = "."

func runWithMockedGolookClient(mockedFunction func()) {

	//ensure that the GolookClient is reset after the function's execution
	defer func(reset client.LookClient) {
		client.GolookClient = reset
	}(client.GolookClient)

	//create a mock client
	client.GolookClient = &MockGolookClient{
		visitDoPostFile: false,
		visitDoPutFiles: false,
		visitDoGetFiles: false,
	}

	mockedFunction()
}

type MockGolookClient struct {
	visitDoPostFile bool
	visitDoPutFiles bool
	visitDoGetFiles bool
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

	t.visitDoPostFile = t.visitDoPostFile || file != nil && file.Name == FILE_NAME
	logrus.WithField("called", t.visitDoPostFile).WithField("file", *file).Info("Test DoPostFile")
	return ""
}

func (t *MockGolookClient) DoPutFiles(files []utils.File) string {
	t.visitDoPutFiles = len(files) > 0
	return ""
}

func (t *MockGolookClient) DoGetFiles() string {
	t.visitDoGetFiles = true
	return ""
}
