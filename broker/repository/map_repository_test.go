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
	"github.com/ottenwbe/golook/broker/models"
	golook "github.com/ottenwbe/golook/broker/runtime/core"
	"github.com/sirupsen/logrus"
	"os"
)

var _ = Describe("The repository implemented with maps", func() {

	var (
		repo *MapRepository
	)

	BeforeEach(func() {
		usePersistence = true
		repo = newMapRepository()
	})

	AfterEach(func() {
		usePersistence = false
		os.RemoveAll(persistenceFile)
	})

	It("does not accept nil systems", func() {
		Expect(repo.StoreSystem("sys", nil)).To(BeFalse())
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
		sysName := sys.Name
		repo.StoreSystem(sysName, sys)
		repo.UpdateFiles(sysName, map[string]*models.File{f.Name: f})

		repo2 := newMapRepository()

		Expect(repo2.systemFiles).To(HaveKey(sysName))
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
		stored := repo.UpdateFiles("unknown", map[string]*models.File{})
		_, found := (repo.systemFiles)["unknown"]
		Expect(stored).To(BeTrue())
		Expect(found).To(BeTrue())
	})

	It("accepts files for valid systems", func() {
		f := newTestFile()
		sys := golook.NewSystem()
		sysName := sys.Name
		repo.StoreSystem(sysName, sys)

		Expect(repo.UpdateFiles(sysName, map[string]*models.File{f.ShortName: f})).To(BeTrue())
	})

	It("should find files that have been stored", func() {
		f := newTestFile()
		sys := golook.NewSystem()
		sysName := sys.Name

		repo.StoreSystem(sysName, sys)
		repo.UpdateFiles(sysName, map[string]*models.File{f.Name: f})

		res := repo.FindSystemAndFiles(f.ShortName)
		Expect(len(res)).To(Equal(1))
		Expect(len(res[sysName])).To(Equal(1))
		Expect(*res[sysName][0]).To(Equal(*f))
	})
})

func newTestFile() *models.File {
	f, err := models.NewFile("map_repository_test.go")
	if err != nil {
		logrus.WithField("file", "map_repository_test.go").Panic("File could not be created in test")
	}
	return f
}
