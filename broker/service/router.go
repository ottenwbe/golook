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
	"github.com/ottenwbe/golook/broker/routing"
	"github.com/ottenwbe/golook/broker/runtime"
)

type Router struct {
	routing.Router
}

var (
	//broadCastRouter's instance is injected during the configuration
	broadCastRouter Router
)

func newRouter(name string, routerTpye routing.RouterType) Router {
	r := Router{routing.NewRouter(name, routerTpye)}

	r.HandlerFunction(SYSTEM_REPORT, handleSystemReport)
	r.HandlerFunction(FILE_QUERY, handleFileQuery)
	r.HandlerFunction(fileReport, handleFileReport)

	routing.ActivateRouter(r)

	newSystemCallbacks.Add("broadcastRouter", r.handleNewSystem)

	return r
}

func (r *Router) handleNewSystem(uuid string, system *runtime.System) {
	r.NewPeer(routing.NewKey(uuid), system.IP)
}
