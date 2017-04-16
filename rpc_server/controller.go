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
package rpc_server

/*
	Common (helper) functions and constants, required by all controllers.
*/

import (
	"net/http"
	"fmt"
	"io"
	"encoding/json"
	"errors"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"

	. "github.com/ottenwbe/golook/models"
	. "github.com/ottenwbe/golook/app"
	. "github.com/ottenwbe/golook/repository"
)

const (
	Nack = "{Nack}"
	Ack  = "{Ack}"

	systemPath = "system"
	filePath   = "file"
)

/////////////////////////////////////
// Helpers for controllers
/////////////////////////////////////

func IsValidRequest(request *http.Request) bool {
	return (request != nil) && (request.Body != nil)
}

func ReturnAck(writer http.ResponseWriter) (int, error) {
	return fmt.Fprint(writer, Ack)
}

func ReturnNackAndLog(writer http.ResponseWriter, errorString string, status int) {
	log.Print(errorString)
	http.Error(writer, errors.New(Nack).Error(), status)
}

func ReturnNackAndLogError(writer http.ResponseWriter, errorString string, err error, status int) {
	log.WithError(err).Print(errorString)
	http.Error(writer, errors.New(Nack).Error(), status)
}

func addSystemAndWriteResult(writer http.ResponseWriter, body io.Reader) {
	system, err := DecodeSystem(body)
	if err != nil {
		fmt.Fprint(writer, Nack)
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

	GoLookRepository.StoreSystem(systemName, system)
	fmt.Fprint(*writer, fmt.Sprintf("{\"id\":\"%s\"}", systemName))
}

func stopOnInvalidRequest(request *http.Request, writer http.ResponseWriter) bool {
	if !IsValidRequest(request) {
		http.Error(writer, errors.New(Nack).Error(), http.StatusBadRequest)
		log.Print("Request rejected: Nil request on server.")
		return true
	}
	return false
}

func marshalFilesAndWriteResult(writer http.ResponseWriter, files map[string]File) {
	if result, marshallErr := json.Marshal(files); marshallErr == nil {
		fmt.Fprintln(writer, string(result))
	} else {
		http.Error(writer, errors.New(Nack).Error(), http.StatusBadRequest)
		log.Printf("Error marshalling file %s", marshallErr)
	}
}

func marshalAndWriteResult(writer http.ResponseWriter, sysFiles map[string]*System) {
	if result, marshallErr := json.Marshal(sysFiles); marshallErr == nil {
		fmt.Fprintln(writer, string(result))
	} else {
		http.Error(writer, errors.New(Nack).Error(), http.StatusBadRequest)
		log.Printf("Error marshalling system/file array: %s", marshallErr)
	}
}

func storeFilesAndWriteResult(system string, files map[string]File, writer http.ResponseWriter) {
	if GoLookRepository.StoreFiles(system, files) {
		ReturnAck(writer)
	} else {
		http.Error(writer, errors.New(Nack).Error(), http.StatusNotFound)
		log.Printf("Files reported from %s could not be stored. System not found.", system)
	}
}

func decodeFilesAndReportSuccess(writer http.ResponseWriter, reader io.Reader, system *string) (map[string]File, bool) {
	if files, err := DecodeFiles(reader); err != nil {
		http.Error(writer, errors.New(Nack).Error(), http.StatusBadRequest)
		log.Printf("Files reported from %s could not be decoded. \n %s", *system, err)
		return nil, false
	} else {
		return files, true
	}
}

func storeFileAndWriteResult(system string, file File, writer http.ResponseWriter) {
	if GoLookRepository.StoreFile(system, file) {
		ReturnAck(writer)
	} else {
		http.Error(writer, errors.New(Nack).Error(), http.StatusNotFound)
		log.Printf("System %s not found while putting a file information to the server.", system)
	}
}

func decodeFileAndReportSuccess(request *http.Request, writer http.ResponseWriter) (File, bool) {
	file, err := DecodeFile(request.Body)
	if err != nil {
		http.Error(writer, errors.New(Nack).Error(), http.StatusBadRequest)
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
