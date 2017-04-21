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
	"fmt"
	"net/http"

	. "github.com/ottenwbe/golook/broker/management"
	. "github.com/ottenwbe/golook/broker/runtime"
)

var (
	queryService  = NewQueryService()
	reportService = NewReportService()
)

// Endpoint: GET /file
func getFiles(writer http.ResponseWriter, request *http.Request) {
	file := extractFileFromPath(request)

	queryService.MakeFileQuery(file)

	returnAck(writer)
}

// Endpoint: GET /info
func getInfo(writer http.ResponseWriter, request *http.Request) {
	fmt.Fprint(writer, EncodeAppInfo(NewAppInfo()))
}

// Endpoint: PUT /file
func putFile(writer http.ResponseWriter, request *http.Request) {

	fileReport, err := extractReport(request)

	if err != nil {
		returnNackAndLogError(writer, "No valid request for: /file", err, http.StatusBadRequest)
	} else {
		reportService.MakeFileReport(fileReport)
		returnAck(writer)
	}
}

// Endpoint: PUT /folder
func putFolder(writer http.ResponseWriter, request *http.Request) {

	folderReport, err := extractReport(request)

	if err != nil {
		returnNackAndLogError(writer, "No valid request for: /folder", err, http.StatusBadRequest)
	} else {
		reportService.MakeFolderReport(folderReport)
		returnAck(writer)
	}
}
