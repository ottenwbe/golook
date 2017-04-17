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

import "encoding/json"

type RequestMessage struct {
	Method  string `json:"encapsulatedMethod"`
	Content string `json:"content"` //TODO: Change to []byte?
}

type ResponseMessage struct {
	Content string `json:"content"` //TODO: Change to []byte?
}

func NewRpcMessage(method string, message interface{}) (*RequestMessage, error) {
	if b, err := json.Marshal(message); err == nil {
		return &RequestMessage{method, string(b)}, nil
	} else {
		return &RequestMessage{"invalid", ""}, err
	}
}

// Get the encapsulated contend. To this end, v is an in/out parameter.
//
// Example:
// m, _ := NewRpcMessage("method", "msg")
// var s string
// m.GetEncapsulated(&s)
//
func (m *RequestMessage) GetEncapsulated(v interface{}) error {
	json.Unmarshal([]byte(m.Content), v)
	return nil
}
