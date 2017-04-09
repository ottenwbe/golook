package control

import "github.com/ottenwbe/golook/utils"

type LookController interface {
	QueryAllSystemsForFile(fileName string) (files map[string]*utils.System, err error)
	QueryReportedFiles() (files []utils.File, err error)
	QueryFiles(systemName string) (files []utils.File, err error)
	ReportFile(filePath string) error
	ReportFolderR(folderPath string) error
	ReportFolder(folderPath string) error
}

type DefaultController struct {
}

func NewController() LookController {
	return DefaultController{}
}
