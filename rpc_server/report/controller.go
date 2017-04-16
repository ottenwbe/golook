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
package report

import (
	"encoding/json"
	. "github.com/ottenwbe/golook/app"
	. "github.com/ottenwbe/golook/file_management"
	. "github.com/ottenwbe/golook/models"
	. "github.com/ottenwbe/golook/rpc_server"
	"net/http"
)

const (
	EP_REPORT = "/report"
)

func init() {
	HttpServer.RegisterFunction(EP_REPORT+"/file", putFile, http.MethodPut)
	HttpServer.RegisterFunction(EP_REPORT+"/folder", putFolder, http.MethodPut)
}

// Endpoint: /report/file
func putFile(writer http.ResponseWriter, request *http.Request) {
	if ! IsValidRequest(request) {
		ReturnNackAndLog(writer, "No valid request for: /report/file", http.StatusBadRequest)
		return
	}

	var fileReport *FileReport
	err := json.NewDecoder(request.Body).Decode(fileReport)

	if err != nil {
		ReturnNackAndLogError(writer, "FileReport could not decoded in controller for: /report/file.", err, http.StatusBadRequest)
		return
	}

	HandleFileReport(fileReport)
	ReturnAck(writer)
}

// Endpoint: /report/folder
func putFolder(writer http.ResponseWriter, request *http.Request) {
	if !IsValidRequest(request) {
		ReturnNackAndLog(writer, "No valid request for: /report/folder", http.StatusBadRequest)
		return
	}

	var fileReport *FileReport
	err := json.NewDecoder(request.Body).Decode(fileReport)
	if err != nil {
		ReturnNackAndLogError(writer, "FileReport could not decoded in controller for: /report/folder.", err, http.StatusBadRequest)
		return
	}

	HandleFolderReport(fileReport)
	ReturnAck(writer)
}
