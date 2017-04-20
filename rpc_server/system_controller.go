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

//
//import (
//	"fmt"
//	"net/http"
//
//	"errors"
//	. "github.com/ottenwbe/golook/app"
//	. "github.com/ottenwbe/golook/repository"
//	log "github.com/sirupsen/logrus"
//	"github.com/gorilla/mux"
//)
//
//const (
//	EP_SYSTEM = "/systems"
//)
//
//func init() {
//	HttpServer.RegisterFunction(EP_SYSTEM, getSystem, http.MethodGet)
//	HttpServer.RegisterFunction(EP_SYSTEM, putSystem, http.MethodPut)
//	HttpServer.RegisterFunction(EP_SYSTEM+"/{system}", getSystemP, http.MethodGet)
//	HttpServer.RegisterFunction(EP_SYSTEM+"/{system}", delSystemP, http.MethodDelete)
//	HttpServer.RegisterFunction(EP_SYSTEM+"/{system}/files", getSystemFiles, http.MethodGet)
//	HttpServer.RegisterFunction(EP_SYSTEM+"/{system}/files/{file}", postFile, http.MethodPost)
//	HttpServer.RegisterFunction(EP_SYSTEM+"/{system}/files", putFiles, http.MethodPut)
//}
//
//// Endpoint: GET /systems
//// Return information about this system
//func getSystem(writer http.ResponseWriter, _ *http.Request) {
//	sys := NewSystem()
//	result := EncodeSystem(sys)
//	fmt.Fprintln(writer, result)
//}
//
//// Endpoint: GET /systems/{system}
//// Return information about a particular system {system}
//func getSystemP(writer http.ResponseWriter, request *http.Request) {
//	system := extractSystemFromPath(request)
//
//	if sys, ok := GoLookRepository.GetSystem(system); ok {
//		writer = marshalSystemAndWritResult(sys, writer)
//	} else {
//		http.Error(writer, errors.New("{}").Error(), http.StatusNotFound)
//		log.Printf("System %s not found: Returning empty Json", system)
//	}
//}
//
//// Endpoint: DELETE /systems/{system}
//// Deletes system {system} on server
//func delSystemP(writer http.ResponseWriter, request *http.Request) {
//	system := extractSystemFromPath(request)
//	GoLookRepository.DelSystem(system)
//	ReturnAck(writer)
//}
//
//
//// Endpoint: GET /systems/{system}/files
//// Get all files of system {system}
//func getSystemFiles(writer http.ResponseWriter, request *http.Request) {
//	params := mux.Vars(request)
//	system := params[systemPath]
//
//	if sys, ok := GoLookRepository.GetSystem(system); ok && sys != nil {
//		marshalFilesAndWriteResult(writer, sys.Files)
//	} else {
//		http.Error(writer, errors.New(Nack).Error(), http.StatusNotFound)
//		log.Printf("Error while receiving files for system %s: system is not registered with sever.", system)
//	}
//}
//
//// Endpoint: PUT /systems/{system}/files
//// Replace all files of system {system}
//func putFiles(writer http.ResponseWriter, request *http.Request) {
//	if stopOnInvalidRequest(request, writer) {
//		return
//	}
//
//	system := extractSystemFromPath(request)
//
//	files, success := decodeFilesAndReportSuccess(writer, request.Body, &system)
//	if !success {
//		return
//	}
//
//	storeFilesAndWriteResult(system, files, writer)
//}
//
//// Endpoint: POST /systems/{system}/files/{file}
//// Add another file to the files stored for a system. Replaces duplicates.
//func postFile(writer http.ResponseWriter, request *http.Request) {
//	if stopOnInvalidRequest(request, writer) {
//		return
//	}
//
//	system := extractSystemFromPath(request)
//
//	file, success := decodeFileAndReportSuccess(request, writer)
//	if !success {
//		return
//	}
//
//	storeFileAndWriteResult(system, file, writer)
//}
//
//// Endpoint: PUT /systems
//// Adds / replaces a system on the server
//func putSystem(writer http.ResponseWriter, request *http.Request) {
//	if stopOnInvalidRequest(request, writer) {
//		return
//	}
//	addSystemAndWriteResult(writer, request.Body)
//}
