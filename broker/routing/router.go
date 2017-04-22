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
	. "github.com/ottenwbe/golook/broker/communication"
	"github.com/ottenwbe/golook/broker/runtime"
	"github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"
)

type Key struct {
	id uuid.UUID
}

func NilKey() Key {
	return Key{
		id: uuid.Nil,
	}
}

func SysKey() Key {
	u, err := uuid.FromString(runtime.GolookSystem.UUID)
	if err != nil {
		log.Error("SysKey() cannot read UUID")
		return NilKey()
	}

	return Key{
		id: u,
	}
}

func NewKeyU(key uuid.UUID) Key {
	return Key{
		id: key,
	}
}

func NewKey(name string) Key {
	return Key{
		id: uuid.NewV5(uuid.Nil, name),
	}
}

func NewKeyN(namespace uuid.UUID, name string) Key {
	return Key{
		id: uuid.NewV5(namespace, name),
	}
}

type HandlerTable map[string]func(params interface{}) interface{}

type RouteTable interface {
	get(key Key) (LookupClient, bool)
	add(key Key, client LookupClient)
	this() LookupClient
}

type DefaultRouteTable struct {
	uplinkClients map[Key]LookupClient
	thisClient    LookupClient
}

func (rt DefaultRouteTable) this() LookupClient {
	return rt.thisClient
}

func (rt DefaultRouteTable) get(key Key) (LookupClient, bool) {
	client, ok := rt.uplinkClients[key]
	return client, ok
}

func (rt DefaultRouteTable) add(key Key, client LookupClient) {
	rt.uplinkClients[key] = client
}

type Router interface {
	Route(key Key, method string, params interface{}) interface{}
	Handle(method string, params interface{}) interface{}
	HandlerFunction(name string, handler func(params interface{}) interface{})
}

type BroadcastRouter struct {
	RouteLayerCallbackClient
	routeTable   DefaultRouteTable //= DefaultRouteTable{}
	routeHandler HandlerTable      //= HandlerTable{}
	name         string
}

func (router BroadcastRouter) HandlerFunction(name string, handler func(params interface{}) interface{}) {
	//TODO: avoid duplicated entries
	router.routeHandler[name] = handler
}

func (router BroadcastRouter) Route(_ Key, method string, message interface{}) (result interface{}) {
	// broadcast to all registered uplink clients
	for _, client := range router.routeTable.uplinkClients {
		// Make the call
		tmpRes, err := client.Call(router.name, method, message)
		if tmpRes != nil && err == nil {
			result = tmpRes
		} else if err != nil {
			log.WithError(err).Error("Route message to client: %s", client)
		}
		log.Debug("Route message to client: %s", client)
	}
	return
}

func (router BroadcastRouter) Handle(method string, message interface{}) interface{} {
	if function, ok := router.routeHandler[method]; ok {
		return function(message)
	} else {
		log.Errorf("Handler for function %s not found", method)
	}
	return nil
}
