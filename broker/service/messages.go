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

import (
	"github.com/ottenwbe/golook/broker/models"
	golook "github.com/ottenwbe/golook/broker/runtime/core"
)

type peerFileReport struct {
	Files      map[string]map[string]*models.File `json:"files"`
	SystemUUID string                             `json:"system"`
}

type peerSystemReport struct {
	SystemUUID string                    `json:"uuid"`
	System     map[string]*golook.System `json:"systems"`
	IsDeletion bool                      `json:"deletion,omitempty"`
	Force      bool                      `json:"force,omitempty"`
}

type peerFileQuery struct {
	SearchString string `json:"search"`
}
