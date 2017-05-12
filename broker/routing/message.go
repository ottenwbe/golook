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

package routing

import (
	golook "github.com/ottenwbe/golook/broker/runtime/core"
	"github.com/ottenwbe/golook/utils"
	log "github.com/sirupsen/logrus"
)

type (
	/*Source represents the source of a message. This source has an id which is typically incremented and a system name*/
	Source struct {
		ID     int    `json:"id"`
		System string `json:"system"`
		//TODO key instead of system
	}

	/*Destination represents the destination in form of a key.*/
	Destination struct {
		Key Key `json:"key"`
	}

	/*Params definition*/
	Params string

	/*RequestMessage represents the generic request from a peer to the next*/
	RequestMessage struct {
		Src    Source      `json:"source"`
		Dst    Destination `json:"destination"`
		Method string      `json:"method"`
		Params Params      `json:"params"`
	}

	/*ResponseMessage represents the generic response to a given request message. The corresponding request can be identified by the RequestSrc field.*/
	ResponseMessage struct {
		RequestSrc  Source `json:"receiver_source"`
		ResponseSrc Source `json:"response_source"`
		Params      Params `json:"params"`
	}
)

/*
NewRequestMessage is the factory function for request messages
*/
func NewRequestMessage(key Key, reqID int, method string, params interface{}) (*RequestMessage, error) {
	p, err := utils.MarshalS(params)
	if err != nil {
		return nil, err
	}
	return &RequestMessage{
		Method: method,
		Params: Params(p),
		Dst:    Destination{Key: key},
		Src:    Source{ID: reqID, System: golook.GolookSystem.UUID},
	}, nil
}

/*
NewResponseMessage is the factory function for response messages
*/
func NewResponseMessage(src Source, params interface{}) (*ResponseMessage, error) {
	p, err := utils.MarshalS(params)
	if err != nil {
		return nil, err
	}
	return &ResponseMessage{
		RequestSrc:  src,
		ResponseSrc: Source{ID: 0, System: golook.GolookSystem.UUID},
		Params:      Params(p),
	}, nil
}

/*
Unmarshal a parameter.
*/
func (p Params) Unmarshal(v interface{}) error {
	log.Debugf("Unmarshalling: %s", string(p))
	err := utils.Unmarshal(string(p), v)
	return err
}

/*
GetEncapsulated returns the encapsulated content of a (rpc) message. To this end, v is an in/out parameter.

Example:
m, _ := NewRequestMessage("method", "msg")
var s string
m.GetEncapsulated(&s)
*/
func (m *RequestMessage) GetEncapsulated(v interface{}) error {
	return utils.Unmarshal(string(m.Params), v)
}
