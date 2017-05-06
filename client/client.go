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

package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/ottenwbe/golook/broker/models"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
)

var (
	/*Host is the url of the golook server.*/
	Host string
)

const (
	/*golookAPIVersion describes the currently supported and implemented http api version*/
	golookAPIVersion = "/v1"
	/*systemEndpoint is the base url for retrieving system information*/
	systemEndpoint = golookAPIVersion + "/system"
	/*fileEndpoint is the base url for querying and monitoring files*/
	fileEndpoint = golookAPIVersion + "/file"
	/*infoEndpoint is the url for basic information about the application*/
	infoEndpoint = "/info"
	/*hTTPApiEndpoint is the url for information about the api*/
	httpAPIEndpoint = "/api"
	/*configEndpoint is the url for querying the configuration*/
	configEndpoint = golookAPIVersion + "/config"
	/*logEndpoint is the url for querying the log*/
	logEndpoint = "/log"
)

/*
GetFiles calls the endpoint /v_/file/{searchString} to find all files similar to searchString on the golook server
*/
func GetFiles(searchString string) (string, error) {
	return get(fmt.Sprintf("%s/%s", fileEndpoint, searchString))
}

/*
ReportFiles will make a file report to the server. In turn, the server will make the given file/folder findable.
*/
func ReportFiles(report models.FileReport) (string, error) {

	log.Debugf("Report for: %s%s", Host, fileEndpoint)

	c := http.Client{}

	b, err := json.Marshal(report)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("%s%s", Host, fileEndpoint), bytes.NewBuffer(b))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.Do(req)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("Error retrieving result from server: %d %s", resp.StatusCode, resp.Status)
	}

	b, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(b), nil
}

/*
GetSystem queries the golook server for system information.
*/
func GetSystem() (string, error) {
	return get(systemEndpoint)
}

/*
GetInfo queries the golook server for general information about the server
*/
func GetInfo() (string, error) {
	return get(infoEndpoint)
}

/*
GetConfig queries the golook server its configuration
*/
func GetConfig() (string, error) {
	return get(configEndpoint)
}

/*
GetAPI queries the golook server for details about the api
*/
func GetAPI() (string, error) {
	return get(httpAPIEndpoint)
}

/*
GetLog queries the golook server for details about its log
*/
func GetLog() (string, error) {
	return get(logEndpoint)
}

func get(endpoint string) (string, error) {
	resp, err := http.Get(fmt.Sprintf("%s%s", Host, endpoint))
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("Error retrieving result from server: %d %s", resp.StatusCode, resp.Status)
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(b), nil
}
