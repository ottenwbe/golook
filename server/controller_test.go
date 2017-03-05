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
package server

import (
	"bytes"
	"encoding/json"
	"github.com/gorilla/mux"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"syscall"
	"testing"
	"time"
)

const systemName = "testSystem"

//Verify that "/" exists and returns the correct status code
func TestHome(t *testing.T) {
	makeTestRequest(
		t,
		"GET",
		"/",
		nil,
		"/",
		home,
		http.StatusOK,
		ack,
	)
}

//Verify that a system can be put to the server, retrieved, and finally be deleted
func TestSystemLifeCycle(t *testing.T) {
	putTestSystem(t, systemName, true)
	defer delTestSystem(t, systemName)

	getTestSystem(t, systemName)
}

func TestSystemLifeCycleWithRandomUuid(t *testing.T) {
	putTestSystem(t, systemName, false)
	//TODO: cleanup
	//defer delTestSystem(t, systemName)
}

//Verify that a file can be created, queried, and finally be deleted on a server, iff a system exists
func TestFileLifeCycle(t *testing.T) {
	putTestSystem(t, systemName, true)
	defer delTestSystem(t, systemName)

	f := createTestFile()
	// create file
	puttTestFile(t, f, systemName)
	// check if file can be retrieved
	getTestFile(t, systemName, f.Name, f.Name)
	// delete file
	putEmptyTestFiles(t, systemName, []byte(`[]`))
	// check if file has been deleted
	getTestFile(t, systemName, f.Name, "")

}

func TestGetNonExistingSystem(t *testing.T) {
	getNonExistingTestSystem(t, systemName)
}

func TestGetFileForNonExistingSystem(t *testing.T) {
	getTestFileForNotExistingSystem(t, systemName, "afile.txt")
}

func TestPutFilesWithWrongJson(t *testing.T) {
	putTestSystem(t, systemName, true)
	defer delTestSystem(t, systemName)

	putFaultyTestFiles(t, systemName, []byte(`[`))
}

func TestPutFileWithWrongJson(t *testing.T) {

	putTestSystem(t, systemName, true)
	defer delTestSystem(t, systemName)

	jsonStr := []byte(`{"file":"afile.txt`)

	makeTestRequest(
		t,
		"Put",
		"/systems/"+systemName+"/files/afile.txt",
		jsonStr,
		"/systems/{system}/files/{file}",
		putFile,
		http.StatusBadRequest,
		nack,
	)
}

func TestGetFiles(t *testing.T) {
	putTestSystem(t, systemName, true)
	defer delTestSystem(t, systemName)

	f := createTestFile()
	// create file
	puttTestFile(t, f, systemName)

	getFiles(t, f.Name, f.Name)

}

/////////////////////////
// Helper Functions - Requests
/////////////////////////

func createTestFile() *File {
	const filename = "controller.go"
	f := &File{}
	fi, err := os.Stat(filename)

	if err != nil {
		log.Fatalf("Test file (%s) could not be opened", filename)
	}

	var stat = fi.Sys().(*syscall.Stat_t)
	f.Accessed = time.Unix(stat.Atim.Sec, stat.Atim.Nsec)
	f.Created = time.Unix(stat.Ctim.Sec, stat.Ctim.Nsec)
	f.Modified = time.Unix(stat.Mtim.Sec, stat.Mtim.Nsec)
	f.Name = filename

	return f
}

func putTestSystem(t *testing.T, name string, useNameAsUuid bool) string {

	var jsonStr []byte

	if useNameAsUuid {
		jsonStr = []byte(`{"name":"1234","os":"linux","ip":"1.1.1.1","uuid":"` + name + `"}`)
	} else {
		jsonStr = []byte(`{"name":"1234","os":"linux","ip":"1.1.1.1"}`)
	}

	return makeTestRequest(
		t,
		"PUT",
		"/systems/"+name,
		jsonStr,
		"/systems/{system}",
		putSystem,
		http.StatusOK,
		"{\"id\":",
	)
}

func getTestSystem(t *testing.T, name string) {
	makeTestRequest(
		t,
		"GET",
		"/systems/"+name,
		nil,
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
		"/systems/"+name,
		nil,
		"/systems/{system}",
		getSystem,
		http.StatusNotFound,
		"{}",
	)
}

func delTestSystem(t *testing.T, name string) {
	makeTestRequest(
		t,
		"DELETE",
		"/systems/"+name,
		nil,
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
		"/systems/"+systemName+"/files/"+filename,
		nil,
		"/systems/{system}/files/{file}",
		getSystemFile,
		http.StatusOK,
		nack,
	)
}

func getTestFile(t *testing.T, systemName string, filename string, comparisonFilename string) {
	makeTestRequest(
		t,
		"GET",
		"/systems/"+systemName+"/files/"+filename,
		nil,
		"/systems/{system}/files/{file}",
		getSystemFile,
		http.StatusOK,
		comparisonFilename,
	)

}

func getFiles(t *testing.T, filename string, comparisonFilename string) {
	makeTestRequest(
		t,
		"GET",
		"/files/"+filename,
		nil,
		"/files/{file}",
		getFile,
		http.StatusOK,
		comparisonFilename,
	)

}

func putFaultyTestFiles(t *testing.T, systemName string, files []byte) {
	makeTestRequest(
		t,
		"PUT",
		"/systems/"+systemName+"/files",
		files,
		"/systems/{system}/files",
		putFiles,
		http.StatusBadRequest,
		nack,
	)
}

func putEmptyTestFiles(t *testing.T, systemName string, files []byte) {
	makeTestRequest(
		t,
		"PUT",
		"/systems/"+systemName+"/files",
		files,
		"/systems/{system}/files",
		putFiles,
		http.StatusOK,
		"",
	)
}

func puttTestFile(t *testing.T, f *File, systemName string) {
	jsonStr, _ := json.Marshal(f)
	makeTestRequest(
		t,
		"PUT",
		"/systems/"+systemName+"/files/"+f.Name,
		jsonStr,
		"/systems/{system}/files/{file}",
		putFile,
		http.StatusOK,
		"",
	)
}

/////////////////////////
// Helper Functions - General HTTP
/////////////////////////

func makeTestRequest(
	t *testing.T,
	testHTTPMethod string,
	toPath string,
	withJson []byte,
	forPattern string,
	withFunction func(http.ResponseWriter, *http.Request),
	expectedResultStatus int,
	expectedResultString string,
) string {
	req, err := http.NewRequest(testHTTPMethod, toPath, makeBody(withJson))
	if err != nil {
		t.Fatalf("Could not make new test request: %s", err)
	}
	determineTestHeaderFromBody(req, withJson)

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

	return rr.Body.String()
}

func makeBody(jsonBody []byte) io.Reader {
	return bytes.NewBuffer(jsonBody)
}

func determineTestHeaderFromBody(req *http.Request, body []byte) {
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
}
