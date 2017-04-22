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
package models

import (
	"encoding/json"
)

type RequestMessage struct {
	Method  string `json:"encapsulatedMethod"`
	Index   string `json:"index"`
	Content string `json:"content"` //TODO: Change to []byte?
}

type ResponseMessage struct {
	Content string `json:"content"` //TODO: Change to []byte?
	Index   string `json:"index"`
}

func Marshal(message interface{}) (interface{}, error) {
	if b, err := json.Marshal(message); err == nil {
		return b, err
	} else {
		return []byte{}, err
	}
}

func MarshalS(message interface{}) (string, error) {
	b, err := Marshal(message)
	return string(b.([]byte)), err
}

func UnmarshalS(message string, result interface{}) error {
	err := Unmarshal([]byte(message), result)
	return err
}

func Unmarshal(message interface{}, result interface{}) error {
	if err := json.Unmarshal(message.([]byte), result); err != nil {
		return err
	}
	return nil
}

func NewRpcMessage(index string, method string, message interface{}) (*RequestMessage, error) {
	m, err := MarshalS(message)
	if err != nil {
		return nil, err
	}
	return &RequestMessage{Method: method, Index: index, Content: m}, nil
}

/*
Get the encapsulated content of a (rpc) message. To this end, v is an in/out parameter.

Example:
m, _ := NewRpcMessage("method", "msg")
var s string
m.GetEncapsulated(&s)
*/
func (m *RequestMessage) GetEncapsulated(v interface{}) error {
	return UnmarshalS(m.Content, v)
}
