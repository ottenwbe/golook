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
	Host string
)

const (
	/*GolookAPIVersion describes the currently supported and implemented http api version*/
	GolookAPIVersion = "/v1"
	/*SystemEndpoint*/
	SystemEndpoint = GolookAPIVersion + "/system"
	/*FileEndpoint is the base url for querying and monitoring files*/
	FileEndpoint = GolookAPIVersion + "/file"
	/*InfoEndpoint is the url for basic information about the application*/
	InfoEndpoint = "/info"
	/*HTTPApiEndpoint is the url for information about the api*/
	HTTPApiEndpoint = "/api"
	/*ConfigEndpoint is the url for querying the configuration*/
	ConfigEndpoint = GolookAPIVersion + "/config"
	/*LogEndpoint is the url for querying the log*/
	LogEndpoint = "/log"
)

func GetFiles(searchString string) (string, error) {
	return get(fmt.Sprintf("%s/%s", FileEndpoint, searchString))
}

func ReportFiles(report models.FileReport) (string, error) {

	log.Debugf("Report for: %s%s", Host, FileEndpoint)

	c := http.Client{}

	b, err := json.Marshal(report)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("%s%s", Host, FileEndpoint), bytes.NewBuffer(b))
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

func GetSystem() (string, error) {
	return get(SystemEndpoint)
}

func GetInfo() (string, error) {
	return get(InfoEndpoint)
}

func GetConfig() (string, error) {
	return get(ConfigEndpoint)
}

func GetApi() (string, error) {
	return get(HTTPApiEndpoint)
}

func GetLog() (string, error) {
	return get(LogEndpoint)
}

func get(ep string) (string, error) {
	resp, err := http.Get(fmt.Sprintf("%s%s", Host, ep))
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
