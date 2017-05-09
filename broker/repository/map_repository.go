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
	golook "github.com/ottenwbe/golook/broker/runtime/core"
	"github.com/sirupsen/logrus"
	"strings"
	"sync"
)

type systemFilesMap map[string]*SystemFiles

/*
MapRepository is the implementation of a repository that stores files and systems in memory in a map.
*/
type MapRepository struct {
	systemFiles systemFilesMap
	mutex       sync.RWMutex
}

func newMapRepository() *MapRepository {
	return &MapRepository{
		systemFiles: make(systemFilesMap, 0),
		mutex:       sync.RWMutex{},
	}
}

func (repo *MapRepository) StoreSystem(name string, system *golook.System) bool {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()

	sys := repo.getOrCreateSystem(name)
	if system != nil {
		sys.System = system
		return true
	}
	return false
}

func (repo *MapRepository) UpdateFiles(name string, files map[string]*File) bool {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()

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

	sys, ok := (repo.systemFiles)[name]
	if !ok {
		sys = &SystemFiles{}
		(repo.systemFiles)[name] = sys
	}
	return sys
}

func (repo *MapRepository) GetSystems() map[string]*golook.System {
	sys := map[string]*golook.System{}
	for id, s := range repo.systemFiles {
		sys[id] = s.System
	}
	return sys
}

func (repo *MapRepository) GetSystem(systemName string) (sys *golook.System, ok bool) {
	repo.mutex.RLock()
	defer repo.mutex.RUnlock()

	system, found := (repo.systemFiles)[systemName]
	if found {
		sys = system.System
	}
	return sys, found
}

func (repo *MapRepository) DelSystem(systemName string) (result *golook.System) {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()

	if r, ok := repo.systemFiles[systemName]; ok {
		result = r.System
	}

	delete(repo.systemFiles, systemName)
	return result
}

/*
GetFiles returns all files stored for a system with systemName. Returns an empty map if the system cannot bee found.
*/
func (repo *MapRepository) GetFiles(systemName string) map[string]*File {
	repo.mutex.RLock()
	defer repo.mutex.RUnlock()

	if sys, found := (repo.systemFiles)[systemName]; found {
		return sys.Files
	}
	return map[string]*File{}
}

//TODO refactor
func (repo *MapRepository) FindSystemAndFiles(findString string) map[string][]*File {
	repo.mutex.RLock()
	defer repo.mutex.RUnlock()

	result := make(map[string][]*File, 0)
	for sid, system := range repo.systemFiles {
		logrus.Info("MapRepository: search for system %s", system)
		for _, file := range system.Files {
			logrus.Info("MapRepository: compare %s vs %", file.Name, findString)
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
