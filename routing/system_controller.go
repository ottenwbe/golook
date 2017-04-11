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
package routing

import (
	. "github.com/ottenwbe/golook/communication"
	. "github.com/ottenwbe/golook/utils"
)

type SystemController struct {
	uplink string
	system *System
}

func NewSystemController(uplink string) *SystemController {
	return &SystemController{
		uplink: uplink,
		system: NewSystem(),
	}
}

func (sc *SystemController) Connect() {
	GolookClient.DoPutSystem(sc.system)
}

func (sc *SystemController) ConnectWith(uplinkHost string) {
	GolookClient.DoPutSystem(sc.system)
}

func (sc *SystemController) Disconnect(uplinkHost string) {

}

func (sc *SystemController) DisconnectAll() {

}
