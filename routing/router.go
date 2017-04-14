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
	. "github.com/ottenwbe/golook/global"
	"github.com/ottenwbe/golook/utils"
)

type LookRouter interface {
	handleQueryAllSystemsForFile(fileName string) (files map[string]*System, err error)
	handleQueryFiles(systemName string) (files map[string]utils.File, err error)
	handleReportFile(filePath string) error
	handleReportFileR(filePath string) error
	handleReportFolderR(folderPath string) error
	handleReportFolder(folderPath string) error
}
