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
package communication

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"

	. "github.com/ottenwbe/golook/repository"
	. "github.com/ottenwbe/golook/utils"
)

const (
	nack = "{nack}"
	ack  = "{ack}"

	systemPath = "system"
	filePath   = "file"
)

var (
	repository Repository
)

func init() {
	repository = NewRepository()
}

/////////////////////////////////////
// Handler for Endpoints
/////////////////////////////////////

// Endpoint: GET /
func home(writer http.ResponseWriter, _ *http.Request) {
	returnAck(writer)
}

// Endpoint: GET /files/{file}
// Get all systems that have matching files to {file}. In addition return information about matching files.
func getFile(writer http.ResponseWriter, request *http.Request) {
	fileName := extractFileFromPath(request)

	sysFiles := repository.FindSystemAndFiles(fileName)

	marshalAndWriteResult(writer, sysFiles)
}

// Endpoint: GET /systems/{system}/files
// Get all files of system {system}
func getSystemFiles(writer http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	system := params[systemPath]

	if sys, ok := repository.GetSystem(system); ok && sys != nil {
		marshalFilesAndWriteResult(writer, sys.Files)
	} else {
		http.Error(writer, errors.New(nack).Error(), http.StatusNotFound)
		log.Printf("Error while receiving files for system %s: system is not registered with sever.", system)
	}
}

// Endpoint: PUT /systems/{system}/files
// Replace all files of system {system}
func putFiles(writer http.ResponseWriter, request *http.Request) {
	if stopOnInvalidRequest(request, writer) {
		return
	}

	system := extractSystemFromPath(request)

	files, success := decodeFilesAndReportSuccess(writer, request.Body, &system)
	if !success {
		return
	}

	storeFilesAndWriteResult(system, files, writer)
}

// Endpoint: POST /systems/{system}/files/{file}
// Add another file to the files stored for a system. Replaces duplicates.
func postFile(writer http.ResponseWriter, request *http.Request) {
	if stopOnInvalidRequest(request, writer) {
		return
	}

	system := extractSystemFromPath(request)

	file, success := decodeFileAndReportSuccess(request, writer)
	if !success {
		return
	}

	storeFileAndWriteResult(system, file, writer)
}

// Endpoint: GET /systems/{system}
func getSystem(writer http.ResponseWriter, request *http.Request) {
	system := extractSystemFromPath(request)

	if sys, ok := repository.GetSystem(system); ok {
		writer = marshalSystemAndWritResult(sys, writer)
	} else {
		http.Error(writer, errors.New("{}").Error(), http.StatusNotFound)
		log.Printf("System %s not found: Returning empty Json", system)
	}
}

// Endpoint: DELETE /systems/{system}
// Deletes system {system} on server
func delSystem(writer http.ResponseWriter, request *http.Request) {
	system := extractSystemFromPath(request)
	repository.DelSystem(system)
	returnAck(writer)
}

// Endpoint: PUT /systems
// Adds / replaces a system on the server
func putSystem(writer http.ResponseWriter, request *http.Request) {
	if stopOnInvalidRequest(request, writer) {
		return
	}
	addSystemAndWriteResult(writer, request.Body)
}

/////////////////////////////////////
// Helpers for controllers
/////////////////////////////////////

func addSystemAndWriteResult(writer http.ResponseWriter, body io.Reader) {
	system, err := DecodeSystem(body)
	if err != nil {
		fmt.Fprint(writer, nack)
		log.Printf("Error: Post system request has errors: %s", err)
	} else {
		formatAndStoreSystem(&system, &writer)
	}
}

func formatAndStoreSystem(system *System, writer *http.ResponseWriter) {
	var systemName string
	if system.UUID == "" {
		systemName = NewUUID()
	} else {
		systemName = system.UUID
	}

	repository.StoreSystem(systemName, system)
	fmt.Fprint(*writer, fmt.Sprintf("{\"id\":\"%s\"}", systemName))
}

func stopOnInvalidRequest(request *http.Request, writer http.ResponseWriter) bool {
	if !isValidRequest(request) {
		http.Error(writer, errors.New(nack).Error(), http.StatusBadRequest)
		log.Print("Request rejected: Nil request on server.")
		return true
	}
	return false
}

func isValidRequest(request *http.Request) bool {
	return (request != nil) && (request.Body != nil)
}

func returnAck(writer http.ResponseWriter) (int, error) {
	return fmt.Fprint(writer, ack)
}

func marshalFilesAndWriteResult(writer http.ResponseWriter, files map[string]File) {
	if result, marshallErr := json.Marshal(files); marshallErr == nil {
		fmt.Fprintln(writer, string(result))
	} else {
		http.Error(writer, errors.New(nack).Error(), http.StatusBadRequest)
		log.Printf("Error marshalling file %s", marshallErr)
	}
}

func marshalAndWriteResult(writer http.ResponseWriter, sysFiles map[string]*System) {
	if result, marshallErr := json.Marshal(sysFiles); marshallErr == nil {
		fmt.Fprintln(writer, string(result))
	} else {
		http.Error(writer, errors.New(nack).Error(), http.StatusBadRequest)
		log.Printf("Error marshalling system/file array: %s", marshallErr)
	}
}

func marshalSystemAndWritResult(sys *System, writer http.ResponseWriter) http.ResponseWriter {
	str, marshalError := json.Marshal(sys)
	if marshalError != nil {
		http.Error(writer, errors.New("{}").Error(), http.StatusNotFound)
		log.Print("Json could not be marshalled")
	} else {
		fmt.Fprint(writer, string(str))
	}
	return writer
}

func storeFilesAndWriteResult(system string, files []File, writer http.ResponseWriter) {
	if repository.StoreFiles(system, files) {
		returnAck(writer)
	} else {
		http.Error(writer, errors.New(nack).Error(), http.StatusNotFound)
		log.Printf("Files reported from %s could not be stored. System not found.", system)
	}
}

func decodeFilesAndReportSuccess(writer http.ResponseWriter, reader io.Reader, system *string) ([]File, bool) {
	if files, err := DecodeFiles(reader); err != nil {
		http.Error(writer, errors.New(nack).Error(), http.StatusBadRequest)
		log.Printf("Files reported from %s could not be decoded. \n %s", *system, err)
		return nil, false
	} else {
		return files, true
	}
}

func storeFileAndWriteResult(system string, file File, writer http.ResponseWriter) {
	if repository.StoreFile(system, file) {
		returnAck(writer)
	} else {
		http.Error(writer, errors.New(nack).Error(), http.StatusNotFound)
		log.Printf("System %s not found while putting a file information to the server.", system)
	}
}

func decodeFileAndReportSuccess(request *http.Request, writer http.ResponseWriter) (File, bool) {
	file, err := DecodeFile(request.Body)
	if err != nil {
		http.Error(writer, errors.New(nack).Error(), http.StatusBadRequest)
		log.WithError(err).Error("File could not be decoded while putting the file to server")
		return File{}, false
	}
	return file, true
}

func extractSystemFromPath(request *http.Request) string {
	params := mux.Vars(request)
	system := params[systemPath]
	return system
}

func extractFileFromPath(request *http.Request) string {
	params := mux.Vars(request)
	fileName := params[filePath]
	return fileName
}
