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
	. "github.com/ottenwbe/golook/broker/models"

	log "github.com/sirupsen/logrus"
)

type BroadcastRouter struct {
	routeTable   RouteTable   //= DefaultRouteTable{}
	routeHandler HandlerTable //= HandlerTable{}
	name         string
	reqId        uint64
}

func newBroadcastRouter(name string) Router {
	return &BroadcastRouter{
		routeTable:   newDefaultRouteTable(),
		routeHandler: HandlerTable{},
		name:         name,
		reqId:        0,
	}
}

func (router *BroadcastRouter) Name() string {
	return router.name
}

func (router *BroadcastRouter) HandlerFunction(name string, handler func(params interface{}) interface{}) {
	//TODO: avoid duplicated entries
	router.routeHandler[name] = handler
}

func (router *BroadcastRouter) BroadCast(method string, message interface{}) interface{} {

	var (
		responseChannel                   = make(chan *ResponseMessage)
		goRoutineCounter uint16           = 0
		responseCounter  uint16           = 0
		result           *ResponseMessage = nil
	)

	m, err := NewRequestMessage(NilKey(), router.reqId, method, message)
	if err != nil {
		log.WithError(err).Error("Request Message could not be created")
		return nil
	}
	router.reqId += 1

	// broadcast to all registered uplink clients
	for _, client := range router.routeTable.clients() {
		go request(client, m, router.name, responseChannel)
		goRoutineCounter += 1
	}

	for result == nil && responseCounter < goRoutineCounter {
		result = <-responseChannel
		responseCounter += 1
	}

	return result
}

func request(client LookupClient, requestMessage *RequestMessage, routerName string, responseChannel chan *ResponseMessage) {
	log.Debug("Routing message to client: %s", client)

	// Make the call
	tmpResponse, err := client.Call(routerName, requestMessage)
	if tmpResponse != nil && err == nil {
		actualResponse := &ResponseMessage{}
		tmpResponse.GetObject(&actualResponse)
		responseChannel <- actualResponse
	} else {
		log.WithError(err).Error("Error while routing message to client: %s", client)
		responseChannel <- nil
	}
}

func (router *BroadcastRouter) Route(_ Key, method string, message interface{}) (result interface{}) {
	return router.BroadCast(method, message)
}

func (router *BroadcastRouter) NewNeighbor(key Key, neighbor LookupClient) {
	log.WithField("router", router.name).WithField("neighbor", neighbor.Url()).Info("New neighbor.")
	router.routeTable.add(key, neighbor)
}

func (router *BroadcastRouter) Handle(routerName string, msg MsgParams) interface{} {

	var (
		result  *ResponseMessage = nil
		message                  = &RequestMessage{}
	)

	// cast message to RequestMessage from interface and verify it is a valid message
	if err := msg.GetObject(&message); err != nil {
		log.WithError(err).Info("Could not read message while handling in router: %s.", routerName)
		return nil
	}

	// ignore duplicates to ensure an at least once semantic
	if duplicateMap.CheckForDuplicates(message.Src) {
		return nil
	}

	// check if it is a broadcast, i.e., NilKey as routing key,
	// then deliver to upper layer and continue with broadcast
	if NilKey() == message.Dst.Key {
		result = router.deliver(message.Method, message.Params).(*ResponseMessage)
		router.BroadCast("", nil)
	} else /* handle directed message */ {
		//check if it is a directed message, then forward it to a possible better match
		//TODO implement directed message...
	}
	return result
}

func (router *BroadcastRouter) deliver(method string, params interface{}) interface{} {
	if handler, ok := router.routeHandler[method]; ok {
		return handler(params)
	} else {
		log.Errorf("Handler for method %s not found in router %s", method, router.name)
	}
	return nil
}
