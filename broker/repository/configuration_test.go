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
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/spf13/viper"
)

var _ = Describe("The configuration of the repositories", func() {

	It("initializes ALL expected default values.", func() {
		InitConfiguration()

		Expect(RepositoryType(viper.GetInt("repository.typeId"))).To(Equal(mapRepository))
		Expect(viper.GetBool("repository.mapRepository.persistence")).To(Equal(false))
		Expect(viper.GetString("repository.mapRepository.persistenceFile")).To(Equal("mapRepository.json"))
	})

	It("applies ALL expected default values.", func() {
		InitConfiguration()
		ApplyConfiguration()

		Expect(repositoryType).To(Equal(mapRepository))
		Expect(defaultMapRepositoryUsePersistence).To(Equal(false))
		Expect(defaultMapRepositoryPersistenceFile).To(Equal("mapRepository.json"))
	})
})
