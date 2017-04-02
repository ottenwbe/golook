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
package repositories

import (
	. "github.com/ottenwbe/golook/utils"
	"strings"
)

type MapRepository map[string]*System

func (repo *MapRepository) StoreSystem(systemName string, system *System) bool {
	if system != nil {
		(*repo)[systemName] = system
		return true
	}
	return false
}

func (repo *MapRepository) StoreFile(systemName string, file File) bool {
	if sys, ok := (*repo)[systemName]; ok {
		if sys.Files == nil {
			sys.Files = make(map[string]File)
		}
		sys.Files[file.Name] = file
		return true
	}
	return false
}

func (repo *MapRepository) StoreFiles(systemName string, files []File) bool {
	if sys, ok := (*repo)[systemName]; ok {
		sys.Files = make(map[string]File)
		for _, file := range files {
			sys.Files[file.Name] = file
		}
		return true
	}
	return false
}

func (repo *MapRepository) GetSystem(systemName string) (sys *System, ok bool) {
	sys, ok = (*repo)[systemName]
	return
}

func (repo *MapRepository) HasFile(fileName string, systemName string) (*File, bool) {
	if sys, ok := repo.GetSystem(systemName); ok {
		if file, ok := sys.Files[fileName]; ok {
			return &file, true
		}

	}
	return nil, false
}

func (repo *MapRepository) DelSystem(systemName string) {
	delete(*repo, systemName)
}

func (repo *MapRepository) FindSystemAndFiles(findString string) map[string]*System {
	result := make(map[string]*System, 0)
	for sid, system := range *repo {
		for _, file := range system.Files {
			if strings.Contains(file.Name, findString) {
				if _, ok := result[sid]; !ok {
					result[sid] = new(System)
					result[sid].Name = system.Name
					result[sid].Files = make(map[string]File, 0)
				}
				result[sid].Files[file.Name] = file
			}
		}
	}
	return result
}
