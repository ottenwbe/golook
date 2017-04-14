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
package rpc

import (
	"fmt"
	"net/http"

	"encoding/json"
	"errors"
	. "github.com/ottenwbe/golook/global"
	. "github.com/ottenwbe/golook/repository"
	log "github.com/sirupsen/logrus"
)

const (
	EP_INFO   = "/info"
	EP_SYSTEM = "/systems"
)

func init() {
	HttpServer.RegisterFunction(EP_INFO, getInfo, "Get")
	HttpServer.RegisterFunction(EP_SYSTEM, getSystem, "Get")
	HttpServer.RegisterFunction(EP_SYSTEM, putSystem, "PUT")
	HttpServer.RegisterFunction(EP_SYSTEM+"/{system}", getSystemP, "GET")
	HttpServer.RegisterFunction(EP_SYSTEM+"/{system}", delSystemP, "DELETE")
	HttpServer.RegisterFunction(EP_SYSTEM+"/{system}/files", getSystemFiles, "GET")
	HttpServer.RegisterFunction(EP_SYSTEM+"/{system}/files/{file}", postFile, "POST")
	HttpServer.RegisterFunction(EP_SYSTEM+"/{system}/files", putFiles, "PUT")
}

func getInfo(writer http.ResponseWriter, _ *http.Request) {
	info := NewAppInfo()
	result := EncodeAppInfo(info)
	fmt.Fprintln(writer, result)
}

// Endpoint: GET /systems
func getSystem(writer http.ResponseWriter, _ *http.Request) {
	sys := NewSystem()
	result := EncodeSystem(sys)
	fmt.Fprintln(writer, result)
}

// Endpoint: GET /systems/{system}
func getSystemP(writer http.ResponseWriter, request *http.Request) {
	system := extractSystemFromPath(request)

	if sys, ok := GoLookRepository.GetSystem(system); ok {
		writer = marshalSystemAndWritResult(sys, writer)
	} else {
		http.Error(writer, errors.New("{}").Error(), http.StatusNotFound)
		log.Printf("System %s not found: Returning empty Json", system)
	}
}

// Endpoint: DELETE /systems/{system}
// Deletes system {system} on server
func delSystemP(writer http.ResponseWriter, request *http.Request) {
	system := extractSystemFromPath(request)
	GoLookRepository.DelSystem(system)
	returnAck(writer)
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
