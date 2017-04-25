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
BroadcastRouter implements a router which by default broadcasts ALL messages to its peers. This
means that direct message requests (see Route) are internally handled as broadcast.
*/
type BroadcastRouter struct {
	routeTable   RouteTable
	routeHandler HandlerTable
	name         string
	reqId        int
}

func newBroadcastRouter(name string) Router {
	return &BroadcastRouter{
		routeTable:   newDefaultRouteTable(),
		routeHandler: HandlerTable{},
		name:         name,
		reqId:        0,
	}
}

func (router *BroadcastRouter) BroadCast(method string, message interface{}) interface{} {

	m, err := NewRequestMessage(NilKey(), router.nextRequestId(), method, message)
	if err != nil {
		log.WithError(err).Error("Request Message could not be created")
		return nil
	}

	resonse := router.broadcast(m)
	if resonse == nil {
		log.Error("Requests response message could not be created")
		return nil
	}

	return resonse.Params
}

func (router *BroadcastRouter) nextRequestId() int {
	result := router.reqId
	router.reqId += 1
	return result
}

func (router *BroadcastRouter) broadcast(m *RequestMessage) *ResponseMessage {
	var (
		responseChannel                   = make(chan *ResponseMessage)
		goRoutineCounter                  = 0
		responseCounter                   = 0
		result           *ResponseMessage = nil
	)

	// forward message to all registered clients concurrently
	for _, client := range router.routeTable.clients() {
		go forward(client, m, router.name, responseChannel)
		goRoutineCounter += 1
	}

	// wait for the first successful response (result != nil) or until all client requests responded
	for result == nil && responseCounter < goRoutineCounter {
		result = <-responseChannel
		responseCounter += 1
	}
	return result
}

func forward(client com.LookupClient, requestMessage *RequestMessage, routerName string, responseChannel chan *ResponseMessage) {
	log.Infof("Routing message to client: %s", client.Url())

	// Make the call
	tmpResponse, err := client.Call(routerName, requestMessage)
	if tmpResponse != nil && err == nil {
		actualResponse := &ResponseMessage{}
		tmpResponse.GetObject(&actualResponse)
		responseChannel <- actualResponse
	} else {
		log.WithError(err).Errorf("Error while routing message to client: %s", client.Url())
		responseChannel <- nil
	}
}

func (router *BroadcastRouter) Route(_ Key, method string, message interface{}) (result interface{}) {
	return router.BroadCast(method, message)
}

func (router *BroadcastRouter) NewPeer(key Key, neighbor com.LookupClient) {
	log.WithField("router", router.name).WithField("neighbor", neighbor.Url()).Info("New neighbor.")
	router.routeTable.add(key, neighbor)
}

func (router *BroadcastRouter) Handle(routerName string, msg models.MsgParams) interface{} {

	var (
		response = ResponseMessage{}
		tmp      = []RequestMessage{}
	)

	// cast request to RequestMessage from interface and verify it is a valid request
	if err := msg.GetObject(&tmp); err != nil {
		log.WithField("router", router.Name()).WithError(err).Infof("Could not read request while handling message.")
		return nil
	}

	request := tmp[0]

	// ignore duplicates to ensure an at least once semantic
	if duplicateMap.CheckForDuplicates(request.Src) {
		return nil
	}

	// check if it is a broadcast, i.e., NilKey as routing key,
	// then deliver to upper layer and continue with broadcast
	if NilKey() == request.Dst.Key {
		responseParams := router.deliver(request.Method, request.Params)
		response = newResponse(responseParams, &request)
		router.broadcast(&request) //TODO forward result of further broadcast
	} else /* handle directed request */ {
		//check if it is a directed request, then forward it to a possible better match
		//TODO implement directed request...
	}
	return response
}

func newResponse(responseParams interface{}, requestMsg *RequestMessage) (result ResponseMessage) {
	if responseParams == nil {
		//TODO err
		r, _ := NewResponseMessage(requestMsg, "{}")
		result = *r
	} else {
		r, _ := NewResponseMessage(requestMsg, responseParams)
		result = *r
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

func (router *BroadcastRouter) Name() string {
	return router.name
}

func (router *BroadcastRouter) HandlerFunction(name string, handler func(params interface{}) interface{}) {
	log.WithField("router", router.Name()).WithField("callback", name).Info("Router registered new callback.")
	router.routeHandler[name] = handler
}
