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

/*
Router dispatches messages. To this end, it implements an embedded routing.Router.
Router informs the embedded routing.Router about possible candidates for peers by registering to the system service's callbacks.
*/
type Router struct {
	routing.Router
}

func newRouter(name string, routerType routing.RouterType) *Router {
	r := &Router{routing.NewRouter(name, routerType)}
	routing.ActivateRouter(r)
	newSystemCallbacks.Add(name, r.handleNewSystem)
	return r
}

func (r *Router) close() {
	routing.DeactivateRouter(r)
	newSystemCallbacks.Delete(r.Router.Name())
}

func (r *Router) handleNewSystem(uuid string, systems map[string]*runtime.System) {
	for _, s := range systems {
		r.NewPeer(routing.NewKey(s.UUID), s.IP)
	}
}
