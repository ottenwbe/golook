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

type RepositoryType int

const ( // iota is reset to 0
	NO_REPOSITORY  RepositoryType = iota // == 0
	MAP_REPOSITORY RepositoryType = iota // == 1
)

var (
	// value is injected through configuration (see configuration.go)
	repositoryType RepositoryType = MAP_REPOSITORY
)

func NewRepository() Repository {
	var repo Repository

	switch repositoryType {
	case NO_REPOSITORY:
		repo = nil
	case MAP_REPOSITORY:
		tmpRepo := make(MapRepository, 0)
		repo = &tmpRepo
	default:
		repo = nil
	}

	return repo
}

// BEWARE: will panic if the GoLookRepository is not a MapRepository
func AccessMapRepository() *MapRepository {
	return GoLookRepository.(*MapRepository)
}

var GoLookRepository Repository = NewRepository()
