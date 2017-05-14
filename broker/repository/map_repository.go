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

/*
MapRepository is the implementation of a repository that stores files and systems in memory in a map.
*/
type MapRepository struct {
	systemFiles     systemFilesMap
	mutex           sync.RWMutex
	persistenceFile string
	usePersistence  bool
}

/*
Types supporting the map repository
*/
type (
	systemFilesMap map[string]*systemFiles

	storedFile struct {
		file    *models.File `json:"file"`
		Monitor bool         `json:"monitor"`
	}

	/*storedFolders represents the file as it is stored in a map_repository*/
	storedFolders struct {
		Files map[string]*storedFile `json:"files"`
	}

	/*systemFiles is a wrapper around files and systems in a common data structure*/
	systemFiles struct {
		System        *golook.System            `json:"system"`
		StoredFolders map[string]*storedFolders `json:"folders"`
	}
)

/*
Configuration parameters of the map repository
*/
var (
	defaultMapRepositoryPersistenceFile = "mapRepo.json"
	defaultMapRepositoryUsePersistence  = false
)

func newMapRepository() *MapRepository {
	result := &MapRepository{
		systemFiles:     systemFilesMap{},
		mutex:           sync.RWMutex{},
		persistenceFile: defaultMapRepositoryPersistenceFile,
		usePersistence:  defaultMapRepositoryUsePersistence,
	}

	result.loadFromDisk()

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

	repo.tryPersist(sysUUID)

	return false
}

/*
UpdateFiles creates or updates a file entry in the map repository.
*/
func (repo *MapRepository) UpdateFiles(sysUUID string, files map[string]map[string]*models.File, monitor bool) bool {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()

	sys := repo.getOrCreateSystem(sysUUID)
	if sys.StoredFolders == nil {
		sys.StoredFolders = map[string]*storedFolders{}
	}

	for folderName, folder := range files {
		sys.tryCreateFolder(folderName)
		sys.handleFiles(folder, folderName, monitor)
	}

	repo.tryPersist(sysUUID)

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
GetSystems returns at runtime all systems for which files can be stored.
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

	repo.tryPersist(sysUUID)

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
		for folderName, folder := range sys.StoredFolders {
			result[folderName] = map[string]*models.File{}
			for fileName, storedFile := range folder.Files {
				result[folderName][fileName] = storedFile.file
			}
		}
	}

	return result
}

/*
GetMonitoredFiles returns all files stored for a system with a given id 'sysUUID'. Returns an empty map if the system cannot be found.
*/
func (repo *MapRepository) GetMonitoredFiles() map[string]map[string]*models.File {
	repo.mutex.RLock()
	defer repo.mutex.RUnlock()

	result := map[string]map[string]*models.File{}

	if sys, found := (repo.systemFiles)[golook.GolookSystem.UUID]; found {
		for folderName, folder := range sys.StoredFolders {
			result[folderName] = map[string]*models.File{}
			for fileName, storedFile := range folder.Files {
				if storedFile.Monitor {
					result[folderName][fileName] = storedFile.file
				}
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
	for systemName, systemFile := range repo.systemFiles {
		log.Debug("MapRepository: search for systemFile %s", systemFile)
		result = systemFile.findFiles(searchString, result, systemName)
	}
	return result
}

func (repo *MapRepository) loadFromDisk() {
	if repo != nil && repo.usePersistence {
		// load persisted file information from disk
		if f, err := ioutil.ReadFile(defaultMapRepositoryPersistenceFile); err == nil {
			systemFile := repo.getOrCreateSystem(golook.GolookSystem.UUID)
			err := utils.Unmarshal(f, &systemFile.StoredFolders)
			if err != nil {
				log.Error("Files cannot be loaded from file.")
			}
		}
	}
}

func (repo *MapRepository) tryPersist(sysUUID string) {
	if sysUUID == golook.GolookSystem.UUID {
		repo.persist()
	}
}

/*persist stores the files that are reported for the system on which golook is running on disk as json*/
func (repo *MapRepository) persist() error {
	systemFile, found := repo.systemFiles[golook.GolookSystem.UUID]
	if repo.usePersistence == true && found {
		return ioutil.WriteFile(repo.persistenceFile, utils.MarshalBD(systemFile), os.ModePerm)
	}
	return nil
}

func (systemFile *systemFiles) findFiles(searchString string, result map[string][]*models.File, systemName string) map[string][]*models.File {
	for _, folder := range systemFile.StoredFolders {
		for _, file := range folder.Files {
			log.Debug("MapRepository: compare %s vs %", file.file.Name, searchString)
			if strings.Contains(file.file.Name, searchString) {
				if _, found := result[systemName]; !found {
					result[systemName] = make([]*models.File, 0)
				}
				result[systemName] = append(result[systemName], file.file)
			}
		}
	}
	return result
}

func (systemFile *systemFiles) handleFiles(folder map[string]*models.File, folderName string, monitor bool) {
	for _, file := range folder {
		if file.Meta.State == models.Created {
			if s, found := systemFile.StoredFolders[folderName].Files[file.Name]; !found {
				systemFile.StoredFolders[folderName].Files[file.Name] = &storedFile{file, monitor}
			} else {
				s.file = file
			}
		} else {
			if file.Directory {
				delete(systemFile.StoredFolders, folderName)
			} else {
				delete(systemFile.StoredFolders[folderName].Files, file.Name)
			}
		}
	}
}

func (systemFile *systemFiles) tryCreateFolder(folderName string) {
	if _, found := systemFile.StoredFolders[folderName]; !found {
		systemFile.StoredFolders[folderName] = &storedFolders{Files: map[string]*storedFile{}}
	}
}
