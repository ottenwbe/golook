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
	. "github.com/ottenwbe/golook/global"
	. "github.com/ottenwbe/golook/utils"

	log "github.com/sirupsen/logrus"

	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type LookupClient interface {
	DoGetHome() string
	DoPostFile(file *File) string
	DoPutFiles(file []File) string
	DoPostFiles(file []File) string
	DoGetSystem(system string) (*System, error)
	DoPutSystem(system *System) *System
	DoDeleteSystem() string
	DoGetFiles(systemName string) (map[string]File, error)
	DoQuerySystemsAndFiles(fileName string) (systems map[string]*System, err error)
}

type LookupClientData struct {
	serverUrl  string
	systemName string
	c          *http.Client //TODO: check if http rpc is synchronized
}

func (lc *LookupClientData) DoGetHome() string {

	response, err := lc.c.Get(lc.serverUrl)
	if err != nil {
		log.Error(err)
	} else {
		defer response.Body.Close()
		ackResponse, _ := ioutil.ReadAll(response.Body)
		return string(ackResponse)
	}
	return ""
}

func (lc *LookupClientData) DoGetSystem(system string) (*System, error) {

	response, err := lc.c.Get(fmt.Sprintf("%s/systems/%s", lc.serverUrl, system))
	if err != nil {
		log.Error(err)
		return nil, err
	} else {
		defer response.Body.Close()
		s, _ := DecodeSystem(response.Body) //TODO error handling
		return &s, nil
	}
}

func (lc *LookupClientData) DoPutSystem(system *System) *System {

	url := fmt.Sprintf("%s/systems", lc.serverUrl)

	jsonBody, _ := json.Marshal(system)
	request, errRequest := http.NewRequest("PUT", url, bytes.NewBuffer(jsonBody))
	if errRequest == nil {
		request.Header.Set("Content-Type", "application/json")
		response, errResult := lc.c.Do(request)
		if errResult != nil {
			log.Error(errResult)
			return &System{}
		} else {
			defer response.Body.Close()
			s, _ := DecodeSystem(response.Body) //TODO error handling
			return &s
		}
	} else {
		log.Error(errRequest)
		return nil
	}
}

func (lc *LookupClientData) DoDeleteSystem() string {

	url := fmt.Sprintf("%s/systems/%s", lc.serverUrl, lc.systemName)

	request, errRequest := http.NewRequest("DELETE", url, nil)
	if errRequest == nil {
		response, errResult := lc.c.Do(request)
		if errResult != nil {
			log.Error(errResult)
		} else {
			defer response.Body.Close()
			res, _ := ioutil.ReadAll(response.Body)
			return string(res) //TODO error handling
		}
	} else {
		log.Error(errRequest)
	}
	return ""
}

func (lc *LookupClientData) DoPostFile(file *File) string {

	log.WithField("file", file.Name).Debug("DoPostFile")

	var fileName string = file.Name

	url := fmt.Sprintf("%s/systems/%s/files/%s", lc.serverUrl, lc.systemName, fileName)

	fileJson, _ := json.Marshal(file) //TODO error handling
	request, errRequest := http.NewRequest("POST", url, bytes.NewBuffer(fileJson))
	if errRequest == nil {
		request.Header.Set("Content-Type", "application/json")
		response, errResult := lc.c.Do(request)
		if errResult != nil {
			log.Error(errResult)
		} else {
			defer response.Body.Close()
			res, _ := ioutil.ReadAll(response.Body)
			return string(res) //TODO error handling
		}
	} else {
		log.Error(errRequest)
	}
	return ""
}

func (lc *LookupClientData) DoPutFiles(file []File) string {

	url := fmt.Sprintf("%s/systems/%s/files", lc.serverUrl, lc.systemName)

	fileJson, _ := json.Marshal(file) //TODO error handling
	request, errRequest := http.NewRequest("PUT", url, bytes.NewBuffer(fileJson))
	if errRequest == nil {
		request.Header.Set("Content-Type", "application/json")
		response, errResult := lc.c.Do(request)
		if errResult != nil {
			log.Error(errResult)
		} else {
			defer response.Body.Close()
			res, _ := ioutil.ReadAll(response.Body)
			return string(res) //TODO error handling
		}
	} else {
		log.Error(errRequest)
	}
	return ""
}

func (lc *LookupClientData) DoPostFiles(file []File) string {

	log.Infof("DoPostFiles for %d files", len(file))

	c := &http.Client{}

	url := fmt.Sprintf("%s/systems/%s/files", lc.serverUrl, lc.systemName)

	fileJson, _ := json.Marshal(file) //TODO error handling
	request, errRequest := http.NewRequest("POST", url, bytes.NewBuffer(fileJson))
	if errRequest == nil {
		request.Header.Set("Content-Type", "application/json")
		response, errResult := c.Do(request)
		if errResult != nil {
			log.Error(errResult)
		} else {
			defer response.Body.Close()
			res, _ := ioutil.ReadAll(response.Body)
			return string(res) //TODO error handling
		}
	} else {
		log.Error(errRequest)
	}
	return ""
}

func (lc *LookupClientData) DoGetFiles(systemName string) (files map[string]File, err error) {

	var (
		response *http.Response
	)

	url := fmt.Sprintf("%s/systems/%s/files", lc.serverUrl, systemName)

	response, err = lc.c.Get(url)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	files, err = DecodeFiles(response.Body)
	return files, err
}

func (lc *LookupClientData) DoQuerySystemsAndFiles(fileName string) (systems map[string]*System, err error) {
	_ = &http.Client{}
	//TODO...
	return nil, nil
}

func NewLookClient(host string, port int) LookupClient {
	return &LookupClientData{
		serverUrl:  fmt.Sprintf("%s:%d", host, port),
		systemName: "",
		c:          &http.Client{},
	}
}
