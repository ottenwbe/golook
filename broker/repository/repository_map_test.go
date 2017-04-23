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

var _ = Describe("The repository map", func() {

	var (
		repo *MapRepository
	)

	BeforeEach(func() {
		repo = &MapRepository{}
	})

	It("not accept nil systems", func() {
		Expect(repo.StoreSystem("sys", nil)).To(BeFalse())
	})

	It("store valid systems", func() {
		sys := runtime.NewSystem()
		Expect(repo.StoreSystem(sys.Name, sys)).To(BeTrue())
		_, ok := (*repo)[sys.Name]
		Expect(ok).To(BeTrue())
	})

	It("retrieve a system by name", func() {
		sys := runtime.NewSystem()
		Expect(repo.StoreSystem(sys.Name, sys)).To(BeTrue())
		sys, ok := repo.GetSystem(sys.Name)
		Expect(sys).ToNot(BeNil())
		Expect(ok).To(BeTrue())
	})

	It("allow to delete a stored System", func() {
		sys := runtime.NewSystem()
		Expect(repo.StoreSystem(sys.Name, sys)).To(BeTrue())
		repo.DelSystem(sys.Name)
		_, ok := (*repo)[sys.Name]
		Expect(ok).To(BeFalse())
	})

	It("should reject files if no valid System is stored", func() {
		Expect(repo.StoreFiles("unknown", map[string]*models.File{})).To(BeFalse())
	})

	It("should accept files for valid systems", func() {
		f := newTestFile()
		sys := runtime.NewSystem()

		Expect(repo.StoreSystem(sys.Name, sys)).To(BeTrue())
		Expect(repo.StoreFiles(sys.Name, map[string]*models.File{f.ShortName: f})).To(BeTrue())
	})

	It("should find files that have been stored", func() {
		f := newTestFile()
		sys := runtime.NewSystem()

		repo.StoreSystem(sys.Name, sys)
		repo.StoreFiles(sys.Name, map[string]*models.File{f.ShortName: f})

		res := repo.FindSystemAndFiles(f.ShortName)
		Expect(len(res)).To(Equal(1))
		Expect(len(res[sys.Name].Files)).To(Equal(1))
		Expect(*res[sys.Name].Files[f.Name]).To(Equal(*f))
	})
})

func newTestFile() *models.File {
	f, err := models.NewFile("repository_map_test.go")
	if err != nil {
		logrus.WithField("file", "repository_map_test.go").Panic("File could not be created in test")
	}
	return f
}
