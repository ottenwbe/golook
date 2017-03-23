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
	"log"
	"net/http"

	"github.com/ottenwbe/golook/client/config"
	. "github.com/ottenwbe/golook/helper"
)

/*func DoGetome() string {
	c := &http.Client{}

	response, err := c.Get(fmt.Sprintf("%s:%d/", config.Host(), config.ServerPort()))
	if err != nil {
		log.Print(err)
	} else {
		defer response.Body.Close()
		return string(response.Body)
	}
	return ""
}
*/

func DoGetSystem(system string) *System {
	c := &http.Client{}

	response, err := c.Get(fmt.Sprintf("%s:%d/systems/%s", config.Host(), config.ServerPort(), system))
	if err != nil {
		log.Print(err)
		return &System{}
	} else {
		defer response.Body.Close()
		s, _ := DecodeSystem(response.Body) //TODO error handling
		return &s
	}
}

/*
func DoPutSystem(system *System) *System {
	c := &http.Client{}

	url := fmt.Sprintf("%s:%d/systems", config.Host(), config.ServerPort())

	jsonBody, _ := json.Marshal(system)
	b := bytes.NewBuffer(jsonBody)
	request, errRequest := http.NewRequest("PUT", url, b)
	if errRequest == nil {
		response, errResult := c.Do(request)
		if errResult != nil {
			log.Print(errResult)
			return &System{}
		} else {
			defer response.Body.Close()
			s, _ := DecodeSystem(response.Body) //TODO error handling
			return &s
		}
	} else {
		log.Print(errRequest)
		return &System{}
	}
}*/

