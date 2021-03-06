//Copyright 2016-2017 Beate Ottenwälder
//
//Licensed under the Apache License, Version 2.0 (the "License");
//you may not use this File except in compliance with the License.
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
	"github.com/spf13/viper"
)

/*
ApplyConfiguration applies the configuration for all sub-components of the repository.
*/
func ApplyConfiguration() {
	repositoryType = RepositoryType(viper.GetInt("repository.typeId"))

	defaultMapRepositoryPersistenceFile = viper.GetString("repository.mapRepository.persistenceFile")
	defaultMapRepositoryUsePersistence = viper.GetBool("repository.mapRepository.persistence")
}

/*
InitConfiguration initializes the configuration for all sub-components of the repository.
*/
func InitConfiguration() {
	viper.SetDefault("repository.typeId", int(mapRepository))

	viper.SetDefault("repository.mapRepository.persistence", false)
	viper.SetDefault("repository.mapRepository.persistenceFile", "mapRepository.json")
}
