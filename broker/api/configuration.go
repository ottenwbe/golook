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
	golook "github.com/ottenwbe/golook/broker/runtime/core"
	"github.com/spf13/viper"
	"net/http"
)

/*
Constants for the configuration of the API layer
*/
const (
	/*GolookAPIVersion describes the currently supported and implemented http api version*/
	GolookAPIVersion = "/v1"

	//SystemPath = "system"

	/*FilePath in an url*/
	FilePath = "file"

	/*SystemEndpoint is the base url for querying and monitoring files*/
	SystemEndpoint = GolookAPIVersion + "/system"
	/*FileEndpoint is the base url for querying and monitoring files*/
	FileEndpoint = GolookAPIVersion + "/file"
	/*QueryEndpoint is the url for querying files*/
	QueryEndpoint = FileEndpoint + "/{" + FilePath + "}"
	/*InfoEndpoint is the url for basic information about the application*/
	InfoEndpoint = "/info"
	/*HTTPApiEndpoint is the url for information about the api*/
	HTTPApiEndpoint = "/api"
	/*ConfigEndpoint is the url for querying the configuration*/
	ConfigEndpoint = GolookAPIVersion + "/config"
	/*LogEndpoint is the url for querying the log*/
	LogEndpoint = "/log"
)

/*
ApplyConfiguration applies the configuration
*/
func ApplyConfiguration() {
	HTTPServer = golook.NewServer(viper.GetString("api.server.address"), golook.ServerHttp)

	HTTPServer.(*golook.HTTPSever).RegisterFunction(FileEndpoint, putFile, http.MethodPut)
	HTTPServer.(*golook.HTTPSever).RegisterFunction(QueryEndpoint, getFiles, http.MethodGet)
	HTTPServer.(*golook.HTTPSever).RegisterFunction(ConfigEndpoint, getConfiguration, http.MethodGet)
	HTTPServer.(*golook.HTTPSever).RegisterFunction(HTTPApiEndpoint, getAPI, http.MethodGet)
	HTTPServer.(*golook.HTTPSever).RegisterFunction(LogEndpoint, getLog, http.MethodGet)
	HTTPServer.(*golook.HTTPSever).RegisterFunction(SystemEndpoint, getSystem, http.MethodGet)

	if viper.GetBool("api.info") {
		HTTPServer.(*golook.HTTPSever).RegisterFunction(InfoEndpoint, getInfo, http.MethodGet)
	}
}

/*
InitConfiguration initializes the configuration
*/
func InitConfiguration() {
	viper.SetDefault("api.info", true)
	viper.SetDefault("api.server.address", ":8383")
}
