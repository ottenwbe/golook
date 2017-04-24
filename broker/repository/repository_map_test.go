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
	"github.com/ottenwbe/golook/broker/runtime"
	"github.com/sirupsen/logrus"
)

var _ = Describe("The repository implemented with maps", func() {

	var (
		repo *MapRepository
	)

	BeforeEach(func() {
		repo = &MapRepository{}
	})

	It("does not accept nil systems", func() {
		Expect(repo.StoreSystem("sys", nil)).To(BeFalse())
	})

	It("stores valid systems", func() {
		sys := &models.SystemFiles{System: runtime.NewSystem(), Files: nil}

		Expect(repo.StoreSystem(sys.System.Name, sys)).To(BeTrue())
		_, ok := (*repo)[sys.System.Name]
		Expect(ok).To(BeTrue())
	})

	It("can retrieve a system by name", func() {
		sys := &models.SystemFiles{System: runtime.NewSystem(), Files: nil}
		repo.StoreSystem(sys.System.Name, sys)

		sys, ok := repo.GetSystem(sys.System.Name)
		Expect(sys).ToNot(BeNil())
		Expect(ok).To(BeTrue())
	})

	It("allows to delete a stored System", func() {
		sys := &models.SystemFiles{System: runtime.NewSystem(), Files: nil}
		sysName := sys.System.Name
		stored := repo.StoreSystem(sysName, sys)

		repo.DelSystem(sysName)
		_, ok := (*repo)[sysName]
		Expect(stored).To(BeTrue())
		Expect(ok).To(BeFalse())
	})

	It("rejects files if no valid System is stored", func() {
		Expect(repo.StoreFiles("unknown", map[string]*models.File{})).To(BeFalse())
	})

	It("accepts files for valid systems", func() {
		f := newTestFile()
		sys := &models.SystemFiles{System: runtime.NewSystem(), Files: nil}
		sysName := sys.System.Name
		repo.StoreSystem(sysName, sys)

		Expect(repo.StoreFiles(sysName, map[string]*models.File{f.ShortName: f})).To(BeTrue())
	})

	It("should find files that have been stored", func() {
		f := newTestFile()
		sys := &models.SystemFiles{System: runtime.NewSystem(), Files: nil}
		sysName := sys.System.Name

		repo.StoreSystem(sysName, sys)
		repo.StoreFiles(sysName, map[string]*models.File{f.ShortName: f})

		res := repo.FindSystemAndFiles(f.ShortName)
		Expect(len(res)).To(Equal(1))
		Expect(len(res[sysName].Files)).To(Equal(1))
		Expect(*res[sysName].Files[f.Name]).To(Equal(*f))
	})
})

func newTestFile() *models.File {
	f, err := models.NewFile("repository_map_test.go")
	if err != nil {
		logrus.WithField("file", "repository_map_test.go").Panic("File could not be created in test")
	}
	return f
}
