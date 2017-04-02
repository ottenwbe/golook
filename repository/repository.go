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

import . "github.com/ottenwbe/golook/utils"

type Repository interface {
	StoreSystem(systemName string, system *System) bool
	GetSystem(systemName string) (*System, bool)
	DelSystem(systemName string)
	HasFile(fileName string, systemName string) (*File, bool)
	StoreFile(systemName string, file File) bool
	StoreFiles(systemName string, files []File) bool
	FindSystemAndFiles(findString string) map[string]*System
}
