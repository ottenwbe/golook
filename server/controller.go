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
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/ottenwbe/golook/helper"
	"io"
)

const (
	nack = "{nack}"
	ack  = "{ack}"
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
func getFile(writer http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	fileName := params[filePath]

	sysFiles := repository.findSystemAndFiles(fileName)

	MarshalAndWriteResult(writer, sysFiles)
}

// Endpoint: GET /systems/{system}/files/{file}
func getSystemFile(writer http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)

	if sys, ok := repository.getSystem(params[systemPath]); ok && sys != nil {
		MarshalFilesAndWriteResult(writer, sys.Files)
	} else {
		fmt.Fprintln(writer, nack)
		log.Print("System to receive file for cannot be found")
	}
}

// Endpoint: PUT /systems/{system}/files
func putFiles(writer http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	system := params[systemPath]

	files, err := DecodeFiles(request.Body) //TODO: ensure request and bond
	if err != nil {
		http.Error(writer, errors.New(nack).Error(), http.StatusBadRequest)
		log.Printf("Files could not be decoded while putting the file to server %s", err)
	} else if repository.storeFiles(system, files) {
		fmt.Fprintln(writer, ack)
	} else {
		http.Error(writer, errors.New(nack).Error(), http.StatusNotFound)
		log.Printf("Files could not be decoded while putting the file to server %s", err)
	}
}

// Endpoint: PUT /systems/{system}/files/{file}
func putFile(writer http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)

	file, err := DecodeFile(request.Body)
	system := params[systemPath]

	if err != nil {
		http.Error(writer, errors.New(nack).Error(), http.StatusBadRequest)
		log.Printf("File could not be decoded while putting the file to server %s", err)
	} else if repository.storeFile(system, file) {
		returnAck(writer)
	} else {
		http.Error(writer, errors.New(nack).Error(), http.StatusNotFound)
		log.Printf("System (%s) not found while putting a file information to the server.", system)
	}
}

// Endpoint: GET /systems/{system}
func getSystem(writer http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	system := params[systemPath]

	if sys, ok := repository.getSystem(system); ok {
		writer = MarshalSystemAndWritResult(sys, writer)
	} else {
		http.Error(writer, errors.New("{}").Error(), http.StatusNotFound)
		log.Printf("System (%s) not found: Returning empty Json", system)
	}
}

// Endpoint: DELETE /systems/{system}
func delSystem(writer http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	repository.delSystem(params[systemPath])
	returnAck(writer)
}

// Endpoint: PUT /system
func putSystem(writer http.ResponseWriter, request *http.Request) {
	if isValidRequest(request) {
		addSystemAndWriteResult(writer, request.Body) //TODO: writer should not be handed down in function
	} else {
		http.Error(writer, errors.New(nack).Error(), http.StatusBadRequest)
		log.Print("Error: Post system request was empty")
	}
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
		extractSystem(&system, &writer)
	}
}

func extractSystem(system *System, writer *http.ResponseWriter) {
	var systemName string
	var errUuid error = nil
	if system.UUID == "" {
		systemName, errUuid = helper.NewUUID()
	} else {
		systemName = system.UUID
	}
	if errUuid == nil {
		repository.storeSystem(systemName, system)
		fmt.Fprint(*writer, "{\"id\":\""+systemName+"\"}")
	} else {
		http.Error(*writer, errors.New(nack).Error(), http.StatusBadRequest)
		log.Printf("Error: UUID error %s", errUuid)
	}
}

func isValidRequest(request *http.Request) bool {
	return (request != nil) && (request.Body != nil)
}

func returnAck(writer http.ResponseWriter) (int, error) {
	return fmt.Fprint(writer, ack)
}

func MarshalFilesAndWriteResult(writer http.ResponseWriter, files []File) {
	if result, marshallErr := json.Marshal(files); marshallErr == nil {
		fmt.Fprintln(writer, string(result))
	} else {
		http.Error(writer, errors.New(nack).Error(), http.StatusBadRequest)
		log.Printf("Error marshalling file %s", marshallErr)
	}
}

func MarshalAndWriteResult(writer http.ResponseWriter, sysFiles map[string]*System) {
	if result, marshallErr := json.Marshal(sysFiles); marshallErr == nil {
		fmt.Fprintln(writer, string(result))
	} else {
		http.Error(writer, errors.New(nack).Error(), http.StatusBadRequest)
		log.Printf("Error marshalling system/file array: %s", marshallErr)
	}
}

func MarshalSystemAndWritResult(sys *System, writer http.ResponseWriter) http.ResponseWriter {
	str, marshalError := json.Marshal(sys)
	if marshalError != nil {
		http.Error(writer, errors.New("{}").Error(), http.StatusNotFound)
		log.Print("Json could not be marshalled")
	} else {
		fmt.Fprint(writer, string(str))
	}
	return writer
}
