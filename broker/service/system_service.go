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
	"github.com/bamzi/jobrunner"
	. "github.com/ottenwbe/golook/broker/models"
	. "github.com/ottenwbe/golook/broker/routing"
	. "github.com/ottenwbe/golook/broker/runtime"
)

func init() {
	//TODO move to runtime package
	jobrunner.Start() // optional: jobrunner.Start(pool int, concurrent int) (10, 1)
	jobrunner.Schedule("@every 5m0s", RouteJob{})
	jobrunner.Now(RouteJob{})
}

// Job Specific Definition
type RouteJob struct {
}

// RouteJob.Run() will get triggered automatically.
func (r RouteJob) Run() {
	routeSystem(false)
}

func routeSystem(isDeletion bool) {
	s := SystemTransfer{Uuid: GolookSystem.UUID, System: &SystemFiles{System: GolookSystem, Files: nil}, IsDeletion: isDeletion}
	systemIndex.Route(SysKey(), SYSTEM_REPORT, s)
}
