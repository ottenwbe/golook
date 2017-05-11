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

/*
RepositoryType represents the repository type
*/
type RepositoryType int

const ( // iota is reset to 0
	noRepository  RepositoryType = iota // == 0
	mapRepository RepositoryType = iota // == 1
)

var (
	// value is injected through configuration (see configuration.go)
	repositoryType = mapRepository
)

/*
NewRepository is the factory function for repositories.
*/
func NewRepository() Repository {
	var repo Repository

	switch repositoryType {
	case noRepository:
		repo = nil
	case mapRepository:
		repo = newMapRepository()
	default:
		repo = nil
	}

	return repo
}

/*
AccessMapRepository casts the GoLookRepository to a map repository; if it is a map repository. Otherwise, the function will panic.
BEWARE: Function will panic if the GoLookRepository is not a MapRepository.
*/
func AccessMapRepository() *MapRepository {
	return GoLookRepository.(*MapRepository)
}

/*
GoLookRepository the global repository
*/
var GoLookRepository = NewRepository()
