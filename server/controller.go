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
	"github.com/gorilla/mux"
	"github.com/ottenwbe/golook/helper"
	"log"
	"net/http"
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

func home(writer http.ResponseWriter, _ *http.Request) {
	fmt.Fprint(writer, ack)
}

// Endpoint: GET /files/{file}
func getFile(writer http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	findString := params[filePath]

	result := repository.findSystemAndFiles(findString)

	//TODO error handling
	bytes, _ := json.Marshal(result)
	fmt.Fprintln(writer, string(bytes))
}

// Endpoint: GET /systems/{system}/files/{file}
func getSystemFile(writer http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)

	if sys, ok := repository.getSystem(params[systemPath]); ok {
		writer = writeFilesOfASystemAsJson(sys, writer)
	} else {
		fmt.Fprintln(writer, nack)
		log.Print("System to receive file for cannot be found")
	}
}

// Endpoint: PUT /systems/{system}/files
func putFiles(writer http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	system := params[systemPath]

	files, err := DecodeFiles(request.Body)
	if err != nil {
		http.Error(writer, errors.New(nack).Error(), http.StatusBadRequest)
	} else if repository.storeFiles(system, files) {
		fmt.Fprintln(writer, ack)
	} else {
		http.Error(writer, errors.New(nack).Error(), http.StatusNotFound)
	}
}

// Endpoint: PUT /systems/{system}/files/{file}
func putFile(writer http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request) //params

	file, err := DecodeFile(request.Body)
	system := params[systemPath]

	if err != nil {
		http.Error(writer, errors.New(nack).Error(), http.StatusBadRequest)
		log.Printf("Error. File could not be decoded while putting the file to server %s", err)
	} else if repository.storeFile(system, file) {
		fmt.Fprintln(writer, ack)
	} else {
		http.Error(writer, errors.New(nack).Error(), http.StatusNotFound)
		log.Printf("Error. System (%s) not found while putting a file information to the server.", system)
	}
}

// Endpoint: GET /systems/{system}
func getSystem(writer http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	system := params[systemPath]

	if sys, ok := repository.getSystem(system); ok {
		str, marshalError := json.Marshal(sys)
		if marshalError != nil {
			http.Error(writer, errors.New("{}").Error(), http.StatusNotFound)
			log.Print("Json could not be marshalled")
		} else {
			fmt.Fprint(writer, string(str))
		}
	} else {
		http.Error(writer, errors.New("{}").Error(), http.StatusNotFound)
		log.Printf("System (%s) not found: Returning empty Json", system)
	}
}

// Endpoint: DELETE /systems/{system}
func delSystem(writer http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	repository.delSystem(params[systemPath])
	fmt.Fprintln(writer, ack)
}

// Endpoint: PUT /systems/{system}
func putSystem(writer http.ResponseWriter, request *http.Request) {
	if (request != nil) && (request.Body != nil) {
		tryAddSystem(&writer, request) //TODO: writer should not be handed down in function
	} else {
		http.Error(writer, errors.New(nack).Error(), http.StatusBadRequest)
		log.Print("Error: Post system request was empty")
	}
}

/////////////////////////////////////
// Helpers for controllers
/////////////////////////////////////

func tryAddSystem(writer *http.ResponseWriter, request *http.Request) {
	system, err := DecodeSystem(request.Body)
	if err != nil {
		log.Printf("Error: Post system request has errors: %s", err)
		fmt.Fprint(*writer, nack)
	} else {
		extractSystem(&system, writer)
	}
}

func extractSystem(system *System, writer *http.ResponseWriter) {
	var newSystem string
	var errUuid error = nil
	if system.UUID == "" {
		newSystem, errUuid = helper.NewUUID()
	} else {
		newSystem = system.UUID
	}
	if errUuid == nil {
		repository.storeSystem(newSystem, system)
		fmt.Fprint(*writer, "{\"id\":\""+newSystem+"\"}")
	} else {
		log.Printf("Error: UUID error %s", errUuid)
		fmt.Fprint(*writer, nack)
	}
}

func writeFilesOfASystemAsJson(sys *System, writer http.ResponseWriter) http.ResponseWriter {
	if bytes, marshallErr := json.Marshal(sys.Files); marshallErr == nil {
		fmt.Fprintln(writer, string(bytes))
	} else {
		fmt.Fprintln(writer, nack)
		log.Printf("Error marshalling file %s", marshallErr)
	}
	return writer
}
