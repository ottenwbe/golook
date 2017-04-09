package control

import (
	"github.com/ottenwbe/golook/routing"
	"github.com/ottenwbe/golook/utils"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	log "github.com/sirupsen/logrus"
)

func TestClients(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Control Test Suite")
}

func RunWithMockedGolookClient(mockedFunction func()) {
	RunWithMockedGolookClientF(mockedFunction, "", "")
}

func RunWithMockedGolookClientF(mockedFunction func(), fileName string, folderName string) {
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
		fileName:         fileName,
		folderName:       folderName,
	}

	mockedFunction()
}

type MockGolookClient struct {
	visitDoPostFile  bool
	visitDoPutFiles  bool
	visitDoGetFiles  bool
	visitDoPostFiles bool
	fileName         string
	folderName       string
}

func (mock *MockGolookClient) DoPostFiles(file []utils.File) string {
	mock.visitDoPostFiles = true
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
	log.WithField("called", mock.visitDoPostFile).WithField("file", *file).Info("Test DoPostFile")
	mock.visitDoPostFile = mock.visitDoPostFile || file != nil && file.Name == mock.fileName
	return ""
}

func (mock *MockGolookClient) DoPutFiles(files []utils.File) string {
	mock.visitDoPutFiles = len(files) > 0
	return ""
}

func (mock *MockGolookClient) DoGetFiles() ([]utils.File, error) {
	mock.visitDoGetFiles = true
	return []utils.File{}, nil
}
