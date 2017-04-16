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
package file_management

import "github.com/bamzi/jobrunner"

// Schedule Type
type FileMonitor interface {
	Start()  //Start Reporting
	Report() //Report Now
}

type DefaultFileMonitor struct {
	frequency string
}

//Frequency can be of the form: "@every 5min60s"
func NewFileMonitor(frequency string) FileMonitor {
	return &DefaultFileMonitor{
		frequency: frequency,
	}
}

func (s *DefaultFileMonitor) Start() {
	jobrunner.Start() // optional: jobrunner.Start(pool int, concurrent int) (10, 1)
	jobrunner.Schedule(s.frequency, ScheduledJob{})
}

func (s *DefaultFileMonitor) Report() {
	jobrunner.Now(ScheduledJob{})
}

type ScheduledJob struct{}

// ScheduledRouter.Run() will get triggered automatically.
func (ScheduledJob) Run() {
	//TODO: route things
}
