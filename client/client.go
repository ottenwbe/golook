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
	"fmt"
	"net/http"

	log "github.com/sirupsen/logrus"

	"bytes"
	"encoding/json"
	"github.com/ottenwbe/golook/config"
	. "github.com/ottenwbe/golook/utils"
	"io/ioutil"
)

type LookClient struct {
	serverUrl string
	//c http.Client //TODO: check if http client is synchronized
}

var GolookClient *LookClient = NewLookClient()

func (lc *LookClient) DoGetHome() string {
	c := &http.Client{}

	response, err := c.Get(lc.serverUrl)
	if err != nil {
		log.Error(err)
	} else {
		defer response.Body.Close()
		ackResponse, _ := ioutil.ReadAll(response.Body)
		return string(ackResponse)
	}
	return ""
}

func (lc *LookClient) DoGetSystem(system string) (*System, error) {
	c := &http.Client{}

	response, err := c.Get(fmt.Sprintf("%s/systems/%s", lc.serverUrl, system))
	if err != nil {
		log.Error(err)
		return nil, err
	} else {
		defer response.Body.Close()
		s, _ := DecodeSystem(response.Body) //TODO error handling
		return &s, nil
	}
}

func (lc *LookClient) DoPutSystem(system *System) *System {
	c := &http.Client{}

	url := fmt.Sprintf("%s/systems", lc.serverUrl)

	jsonBody, _ := json.Marshal(system)
	b := bytes.NewBuffer(jsonBody)
	request, errRequest := http.NewRequest("PUT", url, b)
	if errRequest == nil {
		response, errResult := c.Do(request)
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

func (lc *LookClient) DoDeleteSystem(systemName string) string {
	c := &http.Client{}

	url := fmt.Sprintf("%s/systems/%s", lc.serverUrl, systemName)

	request, errRequest := http.NewRequest("DELETE", url, nil)
	if errRequest == nil {
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

func NewLookClient() *LookClient {
	return &LookClient{
		serverUrl: fmt.Sprintf("%s:%d", config.Host(), config.ServerPort()),
	}
}
