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

package models

import (
	"os"
	"path/filepath"
	"syscall"
	"time"
)

/*
NewFile takes a file path and returns a meta data to that file, iff it exists.
If the file does not exist, we return anyway a file representation and declare it as 'Removed'.
Any other error will be returned.
*/
func NewFile(filePath string) (f *File, err error) {

	fi, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		f = newRemovedFile(filePath)
		return f, nil
	} else if err != nil {
		return nil, err
	}

	fileName, err := filepath.Abs(filePath)
	if err != nil {
		return nil, err
	}

	f = newValidFile(fileName, filePath, fi.Sys().(*syscall.Win32FileAttributeData), fi.IsDir(), fi.Size())
	return f, nil
}

func newRemovedFile(filePath string) *File {
	return &File{
		Name:      filePath,
		ShortName: filePath,
		Created:   time.Unix(0, 0),
		Modified:  time.Unix(0, 0),
		Accessed:  time.Unix(0, 0),
		Meta:      FileMeta{Removed},
	}
}

func newValidFile(fileName string, filePath string, stat *syscall.Win32FileAttributeData, isDir bool, size int64) *File {

	return &File{
		Name:      fileName,
		ShortName: filepath.Base(filePath),
		Created:   time.Unix(0, stat.LastAccessTime.Nanoseconds()),
		Modified:  time.Unix(0, stat.LastWriteTime.Nanoseconds()),
		Accessed:  time.Unix(0, stat.CreationTime.Nanoseconds()),
		Directory: isDir,
		Size:      size,
		Meta:      FileMeta{Created},
	}
}

