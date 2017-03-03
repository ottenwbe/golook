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

package main

import (
	"bytes"
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"syscall"
	"testing"
	"time"
)

//Verify that home exists and returns the correct status code
func TestHome(t *testing.T) {
	makeTestRequest(
		t,
		"GET",
		"/",
		"/",
		home,
		http.StatusOK,
		"{up}",
	)
}

func TestSystemLifeCycle(t *testing.T) {
	const systemName = "testSystem"
	postTestSystem(t, systemName)
	getTestSystem(t, systemName)
	delTestSystem(t, systemName)
}

func TestFileLifeCycle(t *testing.T) {
	const systemName = "testSystem"
	postTestSystem(t, systemName)

	f := createTestFile()
	// create file
	postTestFile(t, f, systemName)
	// check if file can be retrieved
	getTestFile(t, systemName, f.Name, f.Name)
	// delete file
	postEmptyTestFiles(t, systemName)
	// check if file has been deleted
	getTestFile(t, systemName, f.Name, "")

	delTestSystem(t, systemName)
}

func TestGetNonExistingSystem(t *testing.T) {
	const systemName = "testSystem"
	getNonExistingTestSystem(t, systemName)
}

func TestGetFileForNonExistingSystem(t *testing.T) {
	const systemName = "testSystem"
	getTestFileForNotExistingSystem(t, systemName, "afile.txt")
}

/////////////////////////
// Helper Functions
/////////////////////////

func createTestFile() *File {
	f := &File{}
	fi, _ := os.Stat("controller.go")

	var stat = fi.Sys().(*syscall.Stat_t)
	f.Accessed = time.Unix(stat.Atim.Sec, stat.Atim.Nsec)
	f.Created = time.Unix(stat.Ctim.Sec, stat.Ctim.Nsec)
	f.Modified = time.Unix(stat.Mtim.Sec, stat.Mtim.Nsec)
	f.Name = "controller.go"

	return f
}

func postTestSystem(t *testing.T, name string) {
	var jsonStr = []byte(`{"name":"1234","os":"linux","ip":"1.1.1.1","uuid":"` + name + `"}`)
	req, err := http.NewRequest("POST", "/systems/" + name, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	m := mux.NewRouter()
	m.HandleFunc("/systems/{system}", postSystem)
	m.ServeHTTP(rr, req)

	// Expect status 200
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check if the response body is what we expect.
	expected := "{\"id\":"
	if !strings.Contains(rr.Body.String(), expected) {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func getTestSystem(t *testing.T, name string) {
	makeTestRequest(
		t,
		"GET",
		"/systems/" + name,
		"/systems/{system}",
		getSystem,
		http.StatusOK,
		"\"ip\":\"1.1.1.1\"",
	)
}

func getNonExistingTestSystem(t *testing.T, name string) {
	makeTestRequest(
		t,
		"GET",
		"/systems/" + name,
		"/systems/{system}",
		getSystem,
		http.StatusOK,
		"{}",
	)
}

func delTestSystem(t *testing.T, name string) {

	makeTestRequest(
		t,
		"DELETE",
		"/systems/" + name,
		"/systems/{system}",
		delSystem,
		http.StatusOK,
		"Deleting",
	)
}

func getTestFileForNotExistingSystem(t *testing.T, systemName string, filename string) {
	makeTestRequest(
		t,
		"GET",
		"/systems/" + systemName + "/files/" + filename,
		"/systems/{system}/files/{file}",
		getSystemFile,
		http.StatusOK,
		"{nack}",
	)

}


func getTestFile(t *testing.T, systemName string, filename string, comparisonFilename string) {
	makeTestRequest(
		t,
		"GET",
		"/systems/" + systemName + "/files/" + filename,
		"/systems/{system}/files/{file}",
		getSystemFile,
		http.StatusOK,
		comparisonFilename,
	)

}

func postEmptyTestFiles(t *testing.T, systemName string) {
	jsonStr, _ := json.Marshal("[]")
	req, err := http.NewRequest("POST", "/systems/" + systemName + "/files", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	m := mux.NewRouter()
	m.HandleFunc("/systems/{system}/files", putFiles)
	m.ServeHTTP(rr, req)

	// Expect status 200
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

func postTestFile(t *testing.T, f *File, systemName string) {
	jsonStr, _ := json.Marshal(f)
	req, err := http.NewRequest("POST", "/systems/" + systemName + "/files/" + f.Name, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	m := mux.NewRouter()
	m.HandleFunc("/systems/{system}/files/{file}", putFile)
	m.ServeHTTP(rr, req)

	// Expect status 200
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

func makeTestRequest(
t *testing.T,
testHTTPMethod string,
toPath string,
forPattern string,
withFunction func(http.ResponseWriter, *http.Request),
expectedResultStatus int,
expectedResultString string,
) {
	req, err := http.NewRequest(testHTTPMethod, toPath, nil)
	if err != nil {
		t.Fatalf("Could not make new test request: %s", err)
	}

	log.Printf("Make request to: %s %s to test pattern %s", testHTTPMethod, toPath, forPattern)

	rr := httptest.NewRecorder()
	m := mux.NewRouter()
	m.HandleFunc(forPattern, withFunction)
	m.ServeHTTP(rr, req)

	// Expect status
	if status := rr.Code; status != expectedResultStatus {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, expectedResultStatus)
	}

	if !strings.Contains(rr.Body.String(), expectedResultString) {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expectedResultString)
	}
}
