package control

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/ottenwbe/golook/routing"
	"github.com/ottenwbe/golook/utils"
	"github.com/sirupsen/logrus"
	"testing"
)

func TestClients(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Control Test Suite")
}

const FILE_NAME = "reporting_test.go"
const FOLDER_NAME = "."

func runWithMockedGolookClient(mockedFunction func()) {

	//ensure that the GolookClient is reset after the function's execution
	defer func(reset routing.LookClient) {
		routing.GolookClient = reset
	}(routing.GolookClient)

	//create a mock routing
	routing.GolookClient = &MockGolookClient{
		visitDoPostFile:  false,
		visitDoPutFiles:  false,
		visitDoGetFiles:  false,
		visitDoPostFiles: false,
	}

	mockedFunction()
}

type MockGolookClient struct {
	visitDoPostFile  bool
	visitDoPutFiles  bool
	visitDoGetFiles  bool
	visitDoPostFiles bool
}

func (mock *MockGolookClient) DoPostFiles(file []utils.File) string {
	mock.visitDoPostFiles = true
	return ""
}

func (*MockGolookClient) DoQuerySystemsAndFiles() error {
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

	mock.visitDoPostFile = mock.visitDoPostFile || file != nil && file.Name == FILE_NAME
	logrus.WithField("called", mock.visitDoPostFile).WithField("file", *file).Info("Test DoPostFile")
	return ""
}

func (mock *MockGolookClient) DoPutFiles(files []utils.File) string {
	mock.visitDoPutFiles = len(files) > 0
	return ""
}

func (mock *MockGolookClient) DoGetFiles() string {
	mock.visitDoGetFiles = true
	return ""
}
