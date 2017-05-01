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
	"io"
	"time"
)

//FileState represents the current state of the file, i.e., if it has been recently created or removed
type FileState int

const (
	Created FileState = iota
	Removed FileState = iota
)

type File struct {
	Name      string    `json:"name"`
	ShortName string    `json:"short"`
	Created   time.Time `json:"created"`
	Modified  time.Time `json:"modified"`
	Accessed  time.Time `json:"accessed"`
	Directory bool      `json:"directory"`
	Size      int64     `json:"size"`
	Meta      FileMeta  `json:"meta"`
}

type FileMeta struct {
	State FileState `json:"state"`
}

func UnmarshalFiles(fileReader io.Reader) (map[string]File, error) {
	files := make(map[string]File, 0)
	err := json.NewDecoder(fileReader).Decode(&files)
	return files, err
}

func UnmarshalFile(fileReader io.Reader) (File, error) {
	var file File
	err := json.NewDecoder(fileReader).Decode(&file)
	return file, err
}
