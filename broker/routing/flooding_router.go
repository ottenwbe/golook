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
FloodingRouter implements a router which by default delivers ALL messages to its peers. This
means that direct message requests (see Route) are also flooded.
*/
type FloodingRouter struct {
	routeTable    RouteTable
	routeHandlers HandlerTable
	name          string
	reqId         int
}

func newFloodingRouter(name string) Router {
	return &FloodingRouter{
		routeTable:    newDefaultRouteTable(),
		routeHandlers: HandlerTable{},
		name:          name,
		reqId:         0,
	}
}

func (router *FloodingRouter) BroadCast(method string, message interface{}) interface{} {

	m, err := NewRequestMessage(NilKey(), router.nextRequestId(), method, message)
	if err != nil {
		log.WithError(err).Error("Request Message could not be created")
		return nil
	}

	response := router.flood(m)
	if response == nil {
		return nil
	}

	return response.Params
}

func (router *FloodingRouter) nextRequestId() int {
	result := router.reqId
	router.reqId += 1
	return result
}

func (router *FloodingRouter) flood(m *RequestMessage) *ResponseMessage {
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

func forward(client com.RpcClient, request *RequestMessage, router string, responseChannel chan *ResponseMessage) {
	log.WithField("router", router).Debugf("Routing message to client: %s", client.Url())

	// Make the call
	tmpResponse, err := client.Call(router, request)
	if tmpResponse != nil && err == nil {
		actualResponse := &ResponseMessage{}
		tmpResponse.Unmarshal(actualResponse)
		responseChannel <- actualResponse
	} else {
		log.WithError(err).Errorf("Error while routing message to client: %s", client.Url())
		responseChannel <- nil
	}
}

func (router *FloodingRouter) Route(_ Key, method string, message interface{}) (result interface{}) {
	return router.BroadCast(method, message)
}

func (router *FloodingRouter) NewPeer(key Key, neighbor com.RpcClient) {
	log.WithField("router", router.name).WithField("neighbor", neighbor.Url()).Info("New neighbor.")
	router.routeTable.add(key, neighbor)
}

func (router *FloodingRouter) Handle(routerName string, msg models.EncapsulatedValues) interface{} {

	var (
		response *ResponseMessage
		request  = RequestMessage{}
	)

	// cast request to RequestMessage from interface and verify it is a valid request
	if err := msg.Unmarshal(&request); err != nil {
		log.WithField("router", router.Name()).WithError(err).Infof("Could not read request while handling message.")
		return nil
	}

	// ignore duplicates to ensure an at least once semantic
	if duplicateMap.CheckForDuplicates(request.Src) {
		return nil
	}

	// callback to upper layer
	responseParams := router.deliver(request.Method, request.Params)
	response = newResponse(responseParams, &request)

	// treat every message as flooding, therefore:
	// forward message to all other peers
	floodingResponse := router.flood(&request)

	// choose one result
	response = choose(response, floodingResponse)

	return *response
}

func choose(message *ResponseMessage, contender *ResponseMessage) *ResponseMessage {
	if contender == nil || contender.Params == "{}" || contender.Params == "" {
		return message
	}
	if message == nil || message.Params == "{}" || message.Params == "" {
		return contender
	}
	return message
}

func newResponse(responseParams interface{}, requestMsg *RequestMessage) (result *ResponseMessage) {

	if responseParams == nil {
		// Error ignored on purpose since result is nil anyway
		result, _ = NewResponseMessage(requestMsg, "{}")
	} else {
		result, _ = NewResponseMessage(requestMsg, responseParams)
	}

	return result
}

func (router *FloodingRouter) deliver(method string, params interface{}) interface{} {
	if handler, ok := router.routeHandlers[method]; ok {
		return handler(params)
	} else {
		log.Errorf("Handler for method %s not found in router %s", method, router.name)
	}
	return nil
}

/*
Name of the router
*/
func (router *FloodingRouter) Name() string {
	return router.name
}

/*
HandlerFunction registers a handler with the router. The handler is called when a message for this handler arrives.
*/
func (router *FloodingRouter) HandlerFunction(name string, handler func(params interface{}) interface{}) {
	log.WithField("router", router.Name()).WithField("handler", name).Info("Router registered new callback.")
	router.routeHandlers[name] = handler
}
