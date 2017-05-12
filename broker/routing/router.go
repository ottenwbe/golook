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
	log "github.com/sirupsen/logrus"
)

/*
Router describes a generic router interface
*/
type Router interface {
	com.MessageHandler
	Route(key Key, method string, params interface{}) interface{}
	BroadCast(method string, params interface{}) models.EncapsulatedValues
	AddHandler(name string, handler *Handler)
	NewPeer(key Key, url string)
	DelPeer(key Key)
	Name() string
}

/*
Handler is the type of all request and merge handlers
*/
type Handler struct {
	requestHandler func(params models.EncapsulatedValues) interface{}
	mergeCallback  func(raw1 models.EncapsulatedValues, raw2 models.EncapsulatedValues) interface{}
}

/*
NewHandler is the factory method for 'Handler's
*/
func NewHandler(requestHandler func(params models.EncapsulatedValues) interface{}, mergeCallack func(raw1 models.EncapsulatedValues, raw2 models.EncapsulatedValues) interface{}) *Handler {
	return &Handler{
		requestHandler: requestHandler,
		mergeCallback:  mergeCallack,
	}
}

/*
HandlerTable is the alias for a map of handlers
*/
type HandlerTable map[string]*Handler

/*
RouteTable is the base for routing tables
*/
type RouteTable interface {
	peers() map[Key]com.RPCClient
	get(key Key) (com.RPCClient, bool)
	add(key Key, client com.RPCClient)
	del(key Key)
}

func routerLoggerS(rt Router) *log.Entry {
	return log.WithFields(log.Fields{"router": rt.Name()})
}

func routerLogger(rt Router, method string) *log.Entry {
	return log.WithFields(log.Fields{"router": rt.Name(), "method": method})
}

type defaultRouteTable struct {
	peerClients map[Key]com.RPCClient
}

func newDefaultRouteTable() RouteTable {
	return &defaultRouteTable{
		peerClients: make(map[Key]com.RPCClient, 0),
	}
}

func (rt *defaultRouteTable) get(key Key) (com.RPCClient, bool) {
	client, ok := rt.peerClients[key]
	return client, ok
}

func (rt *defaultRouteTable) add(key Key, client com.RPCClient) {
	rt.peerClients[key] = client
}

func (rt *defaultRouteTable) del(key Key) {
	delete(rt.peerClients, key)
}

func (rt *defaultRouteTable) peers() map[Key]com.RPCClient {
	return rt.peerClients
}
