package repositories

import (
	. "github.com/ottenwbe/golook/helper"
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
		sys.Files = append(sys.Files, file)
		return true
	}
	return false
}

func (repo *MapRepository) StoreFiles(systemName string, files []File) bool {
	if sys, ok := (*repo)[systemName]; ok {
		sys.Files = files
		return true
	}
	return false
}

func (repo *MapRepository) GetSystem(systemName string) (sys *System, ok bool) {
	sys, ok = (*repo)[systemName]
	return
}

func (repo *MapRepository) HasFile(fileName string, systemName string) (*File, error) {
	return nil, nil
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
				}
				result[sid].Files = append(result[sid].Files, file)
			}
		}
	}
	return result
}
