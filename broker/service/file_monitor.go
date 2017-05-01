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
	"github.com/fsnotify/fsnotify"
	log "github.com/sirupsen/logrus"
	"sync"
)

/*
FileMonitor
*/
type FileMonitor struct {
	watcher  *fsnotify.Watcher
	once     sync.Once
	done     chan bool
	reporter func(string)
}

func (fm *FileMonitor) Start() {
	var err error

	fm.watcher, err = fsnotify.NewWatcher()
	if err != nil {
		log.WithError(err).Fatal("Cannot start file monitor")
	}

	fm.once.Do(func() {
		fm.done = make(chan bool)
		go cMonitor(fm)
	})
}

func (fm *FileMonitor) Close() {
	if fm.done != nil {
		fm.done <- true
	}
	if fm.watcher != nil {
		fm.watcher.Close()
	}
}

func cMonitor(fm *FileMonitor) {
	var stop bool = false
	for !stop {
		select {
		case event := <-fm.watcher.Events:
			if event.Name != "" {
				log.Infof("Event %s triggered report", event.String())
				if fm.reporter != nil {
					fm.reporter(event.Name)
				} else {
					log.Error("Not reporting monitored file change, since reporter is nil!")
				}
			}
		case err := <-fm.watcher.Errors:
			log.WithError(err).Error("Error from file watcher")
		case stop = <-fm.done:
			log.WithField("stop", stop).Info("Stopping file monitor")
		}
	}
}

/*
Monitor registers paths to files or folders with the FileMonitor. The FileMonitor can then report changes to the fies,
respectively files in the folders.
*/
func (fm *FileMonitor) Monitor(file string) {
	fm.watcher.Add(file)
}

func (fm *FileMonitor) RemoveMonitored(file string) {
	fm.watcher.Remove(file)
}
