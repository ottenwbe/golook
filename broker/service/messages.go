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
package service

import . "github.com/ottenwbe/golook/broker/models"

type PeerFileReport struct {
	Files   map[string]*File `json:"files"`
	Replace bool             `json:"replace"`
	System  string           `json:"system"`
}

type PeerResponse struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
	Data    []byte `json:"data"`
}

type SystemTransfer struct {
	Uuid       string       `json:"uuid"`
	System     *SystemFiles `json:"system"`
	IsDeletion bool         `json:"deletion"`
}

type PeerFileQuery struct {
	SearchString string `json:"search"`
}
