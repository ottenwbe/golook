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
package communication

import (
	"github.com/osamingo/jsonrpc"
	"github.com/ottenwbe/golook/broker/models"
	"github.com/sirupsen/logrus"
)

type RouterCallbackClient interface {
	Handle(method string, params models.MsgParams) interface{}
}

type RegisteredCallbackClients struct {
	routerCallback map[string]RouterCallbackClient
}

func (r *RegisteredCallbackClients) RegisterClient(name string, client RouterCallbackClient, msg interface{}, response interface{}) {
	r.routerCallback[name] = client
	//TODO: registration has to be parametrized
	jsonrpc.RegisterMethod(name, &RpcHandler{
		index: name,
	}, msg, response)

}

func (r *RegisteredCallbackClients) RemoveClient(name string) {
	delete(r.routerCallback, name)
	//TODO remove from jsonrpc, how?
	//jsonrpc.PurgeMethods()
}

func (r *RegisteredCallbackClients) HasClient(name string) bool {
	_, ok := r.routerCallback[name]
	return ok
}

func toRouteLayer(router string, message models.MsgParams) interface{} {

	logrus.Info("going up to route layer for %s", router)

	if reg, ok := RouterRegistrar.routerCallback[router]; ok && reg != nil {
		return reg.Handle(router, message)
	} else {
		logrus.Error("Method dropped before handing it over to route layer. No router client found.")
		return nil
	}

}

var RouterRegistrar = newRouterRegistrar()

func newRouterRegistrar() *RegisteredCallbackClients {
	return &RegisteredCallbackClients{routerCallback: make(map[string]RouterCallbackClient)}
}
