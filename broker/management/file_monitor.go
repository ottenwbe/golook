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
package management

import (
	"github.com/fsnotify/fsnotify"
	log "github.com/sirupsen/logrus"
	"sync"
)

var (
	watcher      *fsnotify.Watcher
	watchedFiles map[string]bool = make(map[string]bool)
	once                         = &sync.Once{}
)

func StartMonitor() {
	var err error

	watcher, err = fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}

	once.Do(func() { go cMonitor() })

	//Future Work: read all previously monitored files from a persistence layer and start monitoring them
}

func cMonitor() {
	for {
		select {
		case event := <-watcher.Events:
			if event.Name != "" {
				log.Infof("Event %s triggered api", event.String())
				routeFile(event.Name, false)
			}
		case err := <-watcher.Errors:
			log.WithError(err).Error("Error from file watcher")
		}
	}
}

func StopMonitor() {
	if watcher != nil {
		watcher.Close()
	}
}

func AddFileMonitor(file string) {
	watcher.Add(file)
	watchedFiles[file] = true
}

func RemoveFileMonitor(file string) {
	watchedFiles[file] = false
	delete(watchedFiles, file)
	watcher.Remove(file)
}

func init() {
	StartMonitor()
}
