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
package api

import (
	"encoding/json"
	"errors"
	"fmt"
	. "github.com/ottenwbe/golook/broker/management"
	. "github.com/ottenwbe/golook/broker/models"
	. "github.com/ottenwbe/golook/broker/runtime"
	"net/http"
)

func ConfigApi() {
	HttpServer.RegisterFunction("/file/{file}", putFile, http.MethodPut)
	HttpServer.RegisterFunction("/file", getFiles, http.MethodGet)
	HttpServer.RegisterFunction("/folder", putFolder, http.MethodPut)
	HttpServer.RegisterFunction("/info", getInfo, http.MethodGet)
}

// Endpoint: GET /file
func getFiles(writer http.ResponseWriter, request *http.Request) {
	file := extractFileFromPath(request)

	MakeFileQuery(file)

	ReturnAck(writer)
}

// Endpoint: GET /info
func getInfo(writer http.ResponseWriter, request *http.Request) {
	fmt.Fprint(writer, EncodeAppInfo(NewAppInfo()))
}

// Endpoint: PUT /file
func putFile(writer http.ResponseWriter, request *http.Request) {

	fileReport, err := extractReport(request)

	if err != nil {
		ReturnNackAndLogError(writer, "No valid request for: /api/file", err, http.StatusBadRequest)
	} else {
		MakeFileReport(fileReport)
		ReturnAck(writer)
	}
}

// Endpoint: PUT /folder
func putFolder(writer http.ResponseWriter, request *http.Request) {

	folderReport, err := extractReport(request)

	if err != nil {
		ReturnNackAndLogError(writer, "No valid request for: /api/folder", err, http.StatusBadRequest)
	} else {
		MakeFolderReport(folderReport)
		ReturnAck(writer)
	}
}

func extractReport(request *http.Request) (*FileReport, error) {
	if !IsValidRequest(request) {
		return nil, errors.New("No valid request")
	}

	var fileReport *FileReport
	err := json.NewDecoder(request.Body).Decode(fileReport)
	if err != nil {
		return nil, err
	}
	return fileReport, nil
}
