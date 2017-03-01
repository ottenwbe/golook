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
package main

import "time"

type File struct {
	Name     string    `json:"name"`
	Created  time.Time `json:"created"`
	Modified time.Time `json:"modified"`
	Accessed time.Time `json:"accessed"`
}

type System struct {
	Name  string `json:"name"`
	OS    string `json:"os"`
	IP    string `json:"ip"`
	UUID  string `json:"uuid"`
	Files []File `json:"files"`
}

type repository interface {
	findFile(fileName string) ([]System, error)
	get(fileName string, systemName string) (File, error)
	putSystem(system System) error
	putFile(file File) error
}

var systemMap = make(map[string]System, 1)