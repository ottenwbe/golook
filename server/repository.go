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
package server

import "strings"

const (
	systemPath = "system"
	filePath   = "file"
)

type Repository interface {
	storeSystem(fileName string, system *System) bool
	getSystem(systemName string) (*System, bool)
	delSystem(systemName string)
	//findFile(fileName string) ([]System, error)
	hasFile(fileName string, systemName string) (*File, error)
	storeFile(systemName string, file File) bool
	storeFiles(systemName string, files []File) bool
	findSystemAndFiles(findString string) map[string]*System
}

func NewRepository() Repository {
	repo := make(RepositoryMap, 0)
	return &repo
}

type RepositoryMap map[string]*System

func (repo *RepositoryMap) storeSystem(systemName string, system *System) bool {
	if system != nil {
		(*repo)[systemName] = system
		return true
	}
	return false
}

func (repo *RepositoryMap) storeFile(systemName string, file File) bool {
	if sys, ok := (*repo)[systemName]; ok {
		sys.Files = append(sys.Files, file)
		return true
	}
	return false
}

func (repo *RepositoryMap) storeFiles(systemName string, files []File) bool {
	if sys, ok := (*repo)[systemName]; ok {
		sys.Files = files
		return true
	}
	return false
}

func (repo *RepositoryMap) getSystem(systemName string) (sys *System, ok bool) {
	sys, ok = (*repo)[systemName]
	return
}

func (repo *RepositoryMap) hasFile(fileName string, systemName string) (*File, error) {
	return nil, nil
}

func (repo *RepositoryMap) delSystem(systemName string) {
	delete(*repo, systemName)
}

func (repo *RepositoryMap) findSystemAndFiles(findString string) map[string]*System {
	result := make(map[string]*System, 0)
	for sid, system := range *repo {
		for _, file := range system.Files {
			if strings.Contains(file.Name, findString) {
				if _, ok := result[sid]; !ok {
					result[sid] = new(System)
					result[sid].Name = system.Name
				}
				result[sid].Files = append(result[sid].Files, file)
			}
		}
	}
	return result
}
