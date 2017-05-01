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
	com "github.com/ottenwbe/golook/broker/communication"
	"github.com/ottenwbe/golook/broker/models"
)

type Router interface {
	com.MessageHandler
	Route(key Key, method string, params interface{}) interface{}
	BroadCast(method string, params interface{}) interface{}
	HandlerFunction(name string, handler func(params models.EncapsulatedValues) interface{})
	NewPeer(key Key, url string)
	Name() string
}

type HandlerTable map[string]func(params models.EncapsulatedValues) interface{}

type RouteTable interface {
	peers() map[Key]com.RpcClient
	get(key Key) (com.RpcClient, bool)
	add(key Key, client com.RpcClient)
	this() com.RpcServer
}

type DefaultRouteTable struct {
	peerClients map[Key]com.RpcClient
}

func newDefaultRouteTable() RouteTable {
	return &DefaultRouteTable{
		peerClients: make(map[Key]com.RpcClient, 0),
	}
}

func (rt *DefaultRouteTable) this() com.RpcServer {
	return nil
}

func (rt *DefaultRouteTable) get(key Key) (com.RpcClient, bool) {
	client, ok := rt.peerClients[key]
	return client, ok
}

func (rt *DefaultRouteTable) add(key Key, client com.RpcClient) {
	rt.peerClients[key] = client
}

func (rt *DefaultRouteTable) peers() map[Key]com.RpcClient {
	return rt.peerClients
}
