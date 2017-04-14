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
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("The repository factory", func() {
	It("should return nil when no repository is configured", func() {
		repositoryType = NO_REPOSITORY
		Expect(NewRepository()).To(BeNil())

	})
	It("should return nil when a wrong repository is configured", func() {
		repositoryType = RepositoryType(2222)
		Expect(NewRepository()).To(BeNil())

	})
	It("should return a map repository when repositoryType is set to MAP_REPOSITORY", func() {
		repositoryType = MAP_REPOSITORY
		repo := NewRepository()
		Expect(repo).ToNot(BeNil())
		//TODO: Check if type is correct
		//Expect(reflect.TypeOf(repo).Name()).To(Equal(reflect.TypeOf(MapRepository{}).Name()))
	})
})
