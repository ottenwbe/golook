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
	"time"
)

/*
FileState represents the current state of the file, i.e., if it has been recently created or removed
*/
type FileState int

/*
Individual file states with the following semantic:
- Created = the file exists and can be found in the file system (fs)
- Removed = the file does not exist in the fs
*/
const (
	Created FileState = iota // == 1
	Removed FileState = iota // == 2
)

/*
File is the logical representation of a file or folder
*/
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

/*
FileMeta describes meta information about a file
*/
type FileMeta struct {
	State FileState `json:"state" repo:"state"`
}
