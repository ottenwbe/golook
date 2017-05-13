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

package core

/*
General information about this application
*/
const (
	GolookName = "golook"
	Version    = "v0.1.0-dev"
)

/*
AppInfo describes an application, i.e., its version and name
*/
type AppInfo struct {
	App     string  `json:"app"`
	Version string  `json:"version"`
	System  *System `json:"system"`
}

/*
NewAppInfo is the factory method for AppInfo
*/
func NewAppInfo() *AppInfo {
	return &AppInfo{
		App:     GolookName,
		Version: Version,
		System:  NewSystem(),
	}
}
