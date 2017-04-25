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
	. "github.com/ottenwbe/golook/broker/runtime"
	. "github.com/ottenwbe/golook/broker/utils"
)

type Source struct {
	Id     int    `json:"id"`
	System string `json:"system"`
}

type Destination struct {
	Key Key `json:"key"`
}

type RequestMessage struct {
	Src    Source      `json:"source"`
	Dst    Destination `json:"destination"`
	Method string      `json:"method"`
	Params string      `json:"params"` //TODO: Change to []byte?
}

type ResponseMessage struct {
	Src      Source `json:"source"`
	Receiver Source `json:"receiver"`
	Params   string `json:"params"` //TODO: Change to []byte?
}

func NewRequestMessage(key Key, reqId int, method string, params interface{}) (*RequestMessage, error) {
	p, err := MarshalS(params)
	if err != nil {
		return nil, err
	}
	return &RequestMessage{Method: method, Params: p, Dst: Destination{Key: key}, Src: Source{Id: reqId, System: GolookSystem.UUID}}, nil
}

func NewResponseMessage(rm *RequestMessage, params interface{}) (*ResponseMessage, error) {
	p, err := MarshalS(params)
	if err != nil {
		return nil, err
	}
	return &ResponseMessage{Src: rm.Src, Receiver: Source{Id: 0, System: GolookSystem.UUID}, Params: p}, nil
}

/*
GetEncapsulated returns the encapsulated content of a (rpc) message. To this end, v is an in/out parameter.

Example:
m, _ := NewRequestMessage("method", "msg")
var s string
m.GetEncapsulated(&s)
*/
func (m *RequestMessage) GetEncapsulated(v interface{}) error {
	return UnmarshalS(m.Params, v)
}
