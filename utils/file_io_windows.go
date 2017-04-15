// +build windows

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

package utils

import (
	"os"
	"path/filepath"
	"syscall"
	"time"
)

func NewFile(filePath string) (f *File, err error) {
	var fi os.FileInfo

	f = &File{}
	fi, err = os.Stat(filePath)
	if err != nil {
		return
	}

	stat := fi.Sys().(*syscall.Win32FileAttributeData)
	f.Accessed = time.Unix(0, stat.LastAccessTime.Nanoseconds())
	f.Created = time.Unix(0, stat.LastWriteTime.Nanoseconds())
	f.Modified = time.Unix(0, stat.CreationTime.Nanoseconds())
	f.Name = filepath.Base(filePath)
	return
}
