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
)

/*
Repository interface for all repositories
*/
type Repository interface {
	StoreSystem(systemName string, system *golook.System) bool
	GetSystem(systemName string) (*golook.System, bool)
	GetSystems() map[string]*golook.System
	DelSystem(systemName string) *golook.System
	UpdateFiles(systemName string, files map[string]*models.File) bool
	FindSystemAndFiles(findString string) map[string][]*models.File
	GetFiles(systemName string) map[string]*models.File
}
