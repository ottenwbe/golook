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
	. "github.com/ottenwbe/golook/app"
	. "github.com/ottenwbe/golook/file_management"
	"testing"
	"time"
)

func TestMapRepository_StoreAndRetrieveSystem(t *testing.T) {
	sysName := "testSys"
	repo := NewRepository()
	s := newTestSystem(sysName)

	repo.StoreSystem(sysName, s)

	if testSys, ok := repo.GetSystem(sysName); !ok || testSys == nil || testSys.IP != s.IP {
		t.Error("System could not be stored in map.")
	}
}

func TestMapRepository_StoreAndRetrieveFile(t *testing.T) {
	sysName := "testSys"
	fileName := "test.txt"
	repo := NewRepository()

	s := newTestSystem(sysName)
	f := newTestFile(fileName)

	repo.StoreSystem(sysName, s)
	repo.StoreFile(sysName, f)

	if _, ok := repo.HasFile(fileName, sysName); !ok {
		t.Error("System could not be stored in map.")
	}
}

func TestMapRepository_FileNotFound(t *testing.T) {
	repo := NewRepository()
	if _, ok := repo.HasFile("test.test", "test"); ok {
		t.Error("File unxpectedly found. Empty map should not comprise any files.")
	}
}

func TestMapRepository_StoreAndRetrieveFiles(t *testing.T) {
	sysName := "testSys"
	fileName := "test.txt"
	fileName2 := "test2.txt"
	repo := NewRepository()
	s := newTestSystem(sysName)
	f := newTestFile(fileName)
	f2 := newTestFile(fileName2)

	repo.StoreSystem(sysName, s)
	repo.StoreFiles(sysName, map[string]File{f.Name: f, f2.Name: f2})

	if _, ok := repo.HasFile(fileName, sysName); !ok {
		t.Error("File could not be found in map.")
	}

	if _, ok := repo.HasFile(fileName2, sysName); !ok {
		t.Error("File could not be found in map.")
	}
}

func TestMapRepository_FindFilesAndSystems(t *testing.T) {
	sysName := "testSys"
	sysName2 := "testSys2"
	fileName := "test.txt"
	fileName2 := "test2.txt"
	repo := NewRepository()
	s := newTestSystem(sysName)
	s2 := newTestSystem(sysName2)
	f := newTestFile(fileName)
	f2 := newTestFile(fileName2)

	repo.StoreSystem(sysName, s)
	repo.StoreSystem(sysName2, s2)
	repo.StoreFiles(sysName, map[string]File{f.Name: f, f2.Name: f2})
	repo.StoreFiles(sysName2, map[string]File{f2.Name: f2})

	testFilesAndSystems := repo.FindSystemAndFiles(fileName)

	if len(testFilesAndSystems) != 1 {
		t.Errorf("Length of test map is suspicous: Expected 1 got %d.", len(testFilesAndSystems))
	} else if testFile, ok := testFilesAndSystems[sysName2]; ok && testFile.Name == f2.Name {
		t.Errorf("System could not be found in test result: Expected %s got %s.", testFile.Name, f2.Name)
	}
}

func TestMapRepository_DelSystems(t *testing.T) {
	sysName := "testSys"
	repo := NewRepository()
	s := newTestSystem(sysName)

	repo.StoreSystem(sysName, s)
	repo.DelSystem(sysName)

	if _, found := repo.GetSystem(sysName); found {
		t.Error("System was not deleted properly")
	}
}

func TestMapRepository_TryStoreInvalidSystem(t *testing.T) {
	sysName := "testSys"
	repo := NewRepository()
	var s *System = nil

	if repo.StoreSystem(sysName, s) {
		t.Error("Expectation not met that nil system can be stored in MapRepository")
	}
}

func TestMapRepository_TryStoreFileOnNonExistingSystem(t *testing.T) {
	fileName := "f.test"
	repo := NewRepository()
	f := newTestFile(fileName)

	if repo.StoreFile("sys", f) {
		t.Error("Expectation not met that repository does not accept files for non existing systems.")
	}

	if repo.StoreFiles("sys", map[string]File{f.Name: f}) {
		t.Error("Expectation not met that repository does not accept files for non existing systems.")
	}
}

func newTestSystem(sysName string) *System {
	s := &System{
		Name:  sysName,
		OS:    "Linux",
		IP:    "localhost",
		UUID:  "uuid",
		Files: nil,
	}
	return s
}

func newTestFile(fileName string) File {
	f := File{}
	f.Name = fileName
	f.Accessed = time.Now()
	f.Created = time.Now()
	f.Modified = time.Now()
	return f
}
