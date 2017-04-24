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
	"net/http"

	. "github.com/ottenwbe/golook/broker/runtime"
)

/*
	Configuration of the api
*/

const (
	Golook_API_VERSION = "/v1"

	//systemPath = "system"
	FILE_PATH = "file"

	FILE_EP       = Golook_API_VERSION + "/file"
	FILE_QUERY_EP = FILE_EP + "/{" + FILE_PATH + "}"
	FOLDER_EP     = Golook_API_VERSION + "/folder"
	INFO_EP       = Golook_API_VERSION + "/info"
)

func RegisterApi() {
	HttpServer.RegisterFunction(FILE_QUERY_EP, putFile, http.MethodPut)
	HttpServer.RegisterFunction(FILE_EP, getFiles, http.MethodGet)
	HttpServer.RegisterFunction(FOLDER_EP, putFolder, http.MethodPut)
	HttpServer.RegisterFunction(INFO_EP, getInfo, http.MethodGet)
}
