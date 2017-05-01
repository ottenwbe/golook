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
	. "github.com/ottenwbe/golook/broker/models"
	"github.com/ottenwbe/golook/broker/runtime"
	"strings"
)

type MapRepository map[string]*SystemFiles

func (repo *MapRepository) StoreSystem(name string, system *runtime.System) bool {
	sys := repo.getOrCreateSystem(name)
	if system != nil {
		sys.System = system
		return true
	}
	return false
}

func (repo *MapRepository) UpdateFiles(name string, files map[string]*File) bool {
	sys := repo.getOrCreateSystem(name)
	if sys.Files == nil {
		sys.Files = make(map[string]*File, 0)
	}
	for _, file := range files {
		if file.Meta.State == Created {
			sys.Files[file.Name] = file
		} else {
			delete(sys.Files, file.Name)
		}
	}
	return true

}
func (repo *MapRepository) getOrCreateSystem(name string) *SystemFiles {
	sys, ok := (*repo)[name]
	if !ok {
		sys = &SystemFiles{}
		(*repo)[name] = sys
	}
	return sys
}

func (repo *MapRepository) GetSystem(systemName string) (sys *runtime.System, ok bool) {
	system, found := (*repo)[systemName]
	if found {
		sys = system.System
	}
	return sys, found
}

func (repo *MapRepository) DelSystem(systemName string) {
	delete(*repo, systemName)
}

func (repo *MapRepository) GetFiles(systemName string) map[string]*File {
	if sys, found := (*repo)[systemName]; found {
		return sys.Files
	}
	return map[string]*File{}
}

//TODO refactor
func (repo *MapRepository) FindSystemAndFiles(findString string) map[string][]*File {
	result := make(map[string][]*File, 0)
	for sid, system := range *repo {
		for _, file := range system.Files {
			if strings.Contains(file.Name, findString) {
				if _, ok := result[sid]; !ok {
					result[sid] = make([]*File, 0)
				}
				result[sid] = append(result[sid], file)
			}
		}
	}
	return result
}
