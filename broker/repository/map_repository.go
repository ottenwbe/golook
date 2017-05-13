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
	"github.com/ottenwbe/golook/broker/models"
	golook "github.com/ottenwbe/golook/broker/runtime/core"
	"github.com/ottenwbe/golook/utils"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"strings"
	"sync"
)

type (
	systemFilesMap map[string]*systemFiles

	/*storedFiles represents the file as it is stored in a map_repository*/
	storedFiles struct {
		Files map[string]*models.File `json:"files"`
	}

	/*systemFiles is a wrapper around files and systems in a common data structure*/
	systemFiles struct {
		System        *golook.System          `json:"system"`
		storedFolders map[string]*storedFiles `json:"folders"`
		TTL           int                     `json:"ttl"`
	}
)

var (
	defaultMapRepositoryPersistenceFile = "mapRepo.json"
	defaultMapRepositoryUsePersistence  = false
)

/*
MapRepository is the implementation of a repository that stores files and systems in memory in a map.
*/
type MapRepository struct {
	systemFiles     systemFilesMap
	mutex           sync.RWMutex
	persistenceFile string
	usePersistence  bool
}

func newMapRepository() *MapRepository {
	result := &MapRepository{
		systemFiles:     systemFilesMap{},
		mutex:           sync.RWMutex{},
		persistenceFile: defaultMapRepositoryPersistenceFile,
		usePersistence:  defaultMapRepositoryUsePersistence,
	}

	loadFromDisk(result)

	return result
}

/*
StoreSystem adds system information to the map repository.
*/
func (repo *MapRepository) StoreSystem(sysUUID string, system *golook.System) bool {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()

	sys := repo.getOrCreateSystem(sysUUID)
	if system != nil {
		sys.System = system
		return true
	}

	tryPersist(sysUUID, repo)

	return false
}

/*
UpdateFiles creates or updates a file entry in the map repository.
*/
func (repo *MapRepository) UpdateFiles(sysUUID string, files map[string]map[string]*models.File) bool {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()

	sys := repo.getOrCreateSystem(sysUUID)
	if sys.storedFolders == nil {
		sys.storedFolders = map[string]*storedFiles{}
	}

	for folderName, folder := range files {
		if _, found := sys.storedFolders[folderName]; !found {
			sys.storedFolders[folderName] = &storedFiles{Files: map[string]*models.File{}}
		}
		for _, file := range folder {
			if file.Meta.State == models.Created {
				sys.storedFolders[folderName].Files[file.Name] = file
			} else {
				if file.Directory {
					delete(sys.storedFolders, folderName)
				} else {
					delete(sys.storedFolders[folderName].Files, file.Name)
				}
			}
		}
	}

	tryPersist(sysUUID, repo)

	return true

}

func (repo *MapRepository) getOrCreateSystem(sysUUID string) *systemFiles {

	systemFile, found := (repo.systemFiles)[sysUUID]
	if !found {
		systemFile = &systemFiles{}
		(repo.systemFiles)[sysUUID] = systemFile
	}
	return systemFile
}

/*
GetSystems returns all systems for which files can be stored.
*/
func (repo *MapRepository) GetSystems() map[string]*golook.System {
	repo.mutex.RLock()
	defer repo.mutex.RUnlock()

	sys := map[string]*golook.System{}
	for systemName, systemFile := range repo.systemFiles {
		sys[systemName] = systemFile.System
	}
	return sys
}

/*
GetSystem returns a system with a given name. If it is not found, ok=false is returned.
*/
func (repo *MapRepository) GetSystem(sysUUID string) (sys *golook.System, ok bool) {
	repo.mutex.RLock()
	defer repo.mutex.RUnlock()

	system, found := (repo.systemFiles)[sysUUID]
	if found {
		sys = system.System
	}
	return sys, found
}

/*
DelSystem removes files and system information from the map. The system which is to be deleted is identified by the systemName.
The method returns the deleted system.
*/
func (repo *MapRepository) DelSystem(sysUUID string) (result *golook.System) {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()

	if systemFile, found := repo.systemFiles[sysUUID]; found {
		result = systemFile.System
	}

	delete(repo.systemFiles, sysUUID)

	tryPersist(sysUUID, repo)

	return result
}

/*
GetFiles returns all files stored for a system with a given id 'sysUUID'. Returns an empty map if the system cannot be found.
*/
func (repo *MapRepository) GetFiles(sysUUID string) map[string]map[string]*models.File {
	repo.mutex.RLock()
	defer repo.mutex.RUnlock()

	result := map[string]map[string]*models.File{}

	if sys, found := (repo.systemFiles)[sysUUID]; found {
		for folderName, folder := range sys.storedFolders {
			result[folderName] = map[string]*models.File{}
			for fileName, storedFile := range folder.Files {
				result[folderName][fileName] = storedFile
			}
		}
	}

	return result
}

/*
FindSystemAndFiles returns all known files that match 'searchString' and the system on which they can be found.
*/
func (repo *MapRepository) FindSystemAndFiles(searchString string) map[string][]*models.File {
	repo.mutex.RLock()
	defer repo.mutex.RUnlock()

	result := make(map[string][]*models.File, 0)
	for sid, system := range repo.systemFiles {
		log.Debug("MapRepository: search for system %s", system)
		result = findFiles(system, searchString, result, sid)
	}
	return result
}

func findFiles(system *systemFiles, searchString string, result map[string][]*models.File, sid string) map[string][]*models.File {
	for _, folder := range system.storedFolders {
		for _, file := range folder.Files {
			log.Debug("MapRepository: compare %s vs %", file.Name, searchString)
			if strings.Contains(file.Name, searchString) {
				if _, found := result[sid]; !found {
					result[sid] = make([]*models.File, 0)
				}
				result[sid] = append(result[sid], file)
			}
		}
	}
	return result
}

func tryPersist(sysUUID string, repo *MapRepository) {
	if sysUUID == golook.GolookSystem.UUID {
		persist(repo)
	}
}

func persist(repo *MapRepository) error {
	systemFile, found := repo.systemFiles[golook.GolookSystem.UUID]
	if repo.usePersistence == true && found {
		return ioutil.WriteFile(repo.persistenceFile, utils.MarshalBD(systemFile), os.ModePerm)
	}
	return nil
}

func loadFromDisk(result *MapRepository) {
	if result != nil && result.usePersistence {
		// load persisted file information from disk
		if f, err := ioutil.ReadFile(defaultMapRepositoryPersistenceFile); err == nil {
			systemFile := result.getOrCreateSystem(golook.GolookSystem.UUID)
			err := utils.Unmarshal(f, &systemFile.storedFolders)
			if err != nil {
				log.Error("Files cannot be loaded from file.")
			}
		}
	}
}
