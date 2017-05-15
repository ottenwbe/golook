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
	"github.com/ottenwbe/golook/broker/models"
	golook "github.com/ottenwbe/golook/broker/runtime/core"
	"github.com/ottenwbe/golook/utils"
	"github.com/sirupsen/logrus"
	"os"
	"path/filepath"
)

var _ = Describe("The repository implemented with maps", func() {

	var (
		repo *MapRepository
	)

	BeforeEach(func() {
		repo = newMapRepository()
		repo.usePersistence = true
	})

	AfterEach(func() {
		os.RemoveAll(repo.persistenceFile)
	})

	It("does not accept nil systems", func() {
		Expect(repo.StoreSystem("systemFile", nil)).To(BeFalse())
	})

	It("stores valid systems", func() {
		sys := golook.NewSystem()

		Expect(repo.StoreSystem(sys.Name, sys)).To(BeTrue())
		_, ok := (repo.systemFiles)[sys.Name]
		Expect(ok).To(BeTrue())
	})

	It("can retrieve a system by name", func() {
		sys := golook.NewSystem()
		repo.StoreSystem(sys.Name, sys)

		sys, ok := repo.GetSystem(sys.Name)
		Expect(sys).ToNot(BeNil())
		Expect(ok).To(BeTrue())
	})

	It("can be read from disk", func() {
		f := newTestFile()
		sys := golook.NewSystem()
		sysUUID := sys.UUID
		repo.StoreSystem(sysUUID, sys)
		repo.UpdateFiles(sysUUID, makeTestFolder(f), false)

		var persist = true
		utils.Mock(&defaultMapRepositoryUsePersistence, &persist, func() {
			repo2 := newMapRepository()
			Expect(repo2.systemFiles).To(HaveKey(sysUUID))
		})
	})

	It("allows to delete a stored System", func() {
		sys := golook.NewSystem()
		sysName := sys.Name
		stored := repo.StoreSystem(sysName, sys)

		repo.DelSystem(sysName)
		_, ok := (repo.systemFiles)[sysName]
		Expect(stored).To(BeTrue())
		Expect(ok).To(BeFalse())
	})

	It("accepts files if no valid System is stored and creates an entry for that system", func() {
		stored := repo.UpdateFiles("unknown", map[string]map[string]*models.File{}, false)
		_, found := (repo.systemFiles)["unknown"]
		Expect(stored).To(BeTrue())
		Expect(found).To(BeTrue())
	})

	It("accepts files for valid systems", func() {
		f := newTestFile()
		sys := golook.NewSystem()
		sysName := sys.Name
		repo.StoreSystem(sysName, sys)

		Expect(repo.UpdateFiles(sysName, makeTestFolder(f), false)).To(BeTrue())
	})

	It("should find files that have been stored", func() {
		f := newTestFile()
		sys := golook.NewSystem()
		sysName := sys.Name

		repo.StoreSystem(sysName, sys)
		repo.UpdateFiles(sysName, makeTestFolder(f), false)

		res := repo.FindSystemAndFiles(f.ShortName)
		Expect(len(res)).To(Equal(1))
		Expect(len(res[sysName])).To(Equal(1))
		Expect(*res[sysName][0]).To(Equal(*f))
	})
})

func newTestFile() *models.File {
	f, err := models.NewFile("map_repository_test.go")
	if err != nil {
		logrus.WithField("File", "map_repository_test.go").Panic("Files could not be created in test")
	}
	return f
}

func makeTestFolder(f *models.File) map[string]map[string]*models.File {
	return map[string]map[string]*models.File{filepath.Dir(f.Name): {f.Name: f}}
}
