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
	"errors"
	"fmt"
	golook "github.com/ottenwbe/golook/broker/runtime/core"
	"github.com/ottenwbe/golook/broker/service"
	"github.com/ottenwbe/golook/utils"
	"net/http"
)

var (
	configurationService = service.NewConfigurationService()
	fileServices         = service.OpenFileServices(service.BroadcastFiles)
)

//getFiles implements the http endpoint: GET /file
func getFiles(writer http.ResponseWriter, request *http.Request) {
	file := extractFileFromPath(request)

	result, err := fileServices.Query(file)
	if err != nil {
		returnAndLogError(writer, "Cannot process query.", err, http.StatusInternalServerError)
	}

	fmt.Fprint(writer, result)
}

//putFile implements the Endpoint: PUT /file
func putFile(writer http.ResponseWriter, request *http.Request) {

	fileReport, err := extractReport(request)
	if err != nil {
		returnAndLogError(writer, "No valid request for: /file", err, http.StatusBadRequest)
		return
	}

	files, err := fileServices.Report(fileReport)
	if err != nil {
		returnAndLogError(writer, "No valid request for: /file", err, http.StatusBadRequest)
		return
	}

	jsonResult, err := utils.MarshalS(files)
	if err != nil {
		returnAndLogError(writer, "Cannot marshal response.", err, http.StatusInternalServerError)
		return
	}

	fmt.Fprint(writer, jsonResult)

}

//getAPI implements the http Endpoint: GET /api
func getAPI(writer http.ResponseWriter, _ *http.Request) {

	if HTTPServer == nil {
		returnAndLogError(writer, "HTTPServer appears to be down.", nil, http.StatusInternalServerError)
		return
	}

	infoAPI := HTTPServer.Info()

	jsonResponse, err := utils.MarshalS(infoAPI)
	if err != nil {
		returnAndLogError(writer, "Cannot marshal Api information.", err, http.StatusInternalServerError)
		return
	}

	fmt.Fprint(writer, jsonResponse)
}

//getLog implements the Endpoint: GET /log
func getLog(writer http.ResponseWriter, _ *http.Request) {
	err := service.RewriteLog(writer)
	if err != nil {
		returnAndLogError(writer, "Cannot open log file", err, http.StatusInternalServerError)
		return
	}
}

//getSystem implements the Endpoint: GET /system
func getSystem(writer http.ResponseWriter, _ *http.Request) {

	systems := service.GetSystems()

	jsonResponse, err := utils.MarshalS(systems)
	if err != nil {
		returnAndLogError(writer, "Cannot marshal systems.", err, http.StatusInternalServerError)
		return
	}

	fmt.Fprint(writer, jsonResponse)
}

//getConfiguration implements the Endpoint: GET /config
func getConfiguration(writer http.ResponseWriter, _ *http.Request) {

	if configurationService == nil {
		returnAndLogError(writer, "ConfigurationService is nil.", errors.New("ConfigurationService is nil."), http.StatusInternalServerError)
		return
	}

	configurations := configurationService.GetConfiguration()

	jsonResponse, err := utils.MarshalS(configurations)
	if err != nil {
		returnAndLogError(writer, "Cannot marshal configuration.", err, http.StatusInternalServerError)
		return
	}

	fmt.Fprint(writer, jsonResponse)
}

//getInfo implements the Endpoint: GET /info
func getInfo(writer http.ResponseWriter, _ *http.Request) {
	jsonResponse := utils.MarshalSD(golook.NewAppInfo())
	fmt.Fprint(writer, jsonResponse)
}
