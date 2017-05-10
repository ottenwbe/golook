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

/*
BroadCastRouter implements a router which by default delivers ALL messages to its peerClients. This
means that direct message requests (see Route) are also flooded.
*/
type BroadCastRouter struct {
	name          string
	routeTable    RouteTable
	routeHandlers HandlerTable
	reqId         int
	duplicateMap  *duplicateMap
}

func newBroadcastRouter(name string) Router {
	return &BroadCastRouter{
		routeTable:    newDefaultRouteTable(),
		routeHandlers: HandlerTable{},
		name:          name,
		reqId:         0,
		duplicateMap:  newDuplicateMap(),
	}
}

/*
BroadCast a message
*/
func (router *BroadCastRouter) BroadCast(method string, message interface{}) models.EncapsulatedValues {

	m, err := NewRequestMessage(NilKey(), router.nextRequestId(), method, message)
	if err != nil {
		routerLogger(router, method).WithError(err).Error("Cannot create Request Message.")
		return nil
	}

	routerLogger(router, method).
		WithField("request_id", router.reqId).
		WithField("msg_id", m.Src.Id).
		Info("New request message will be sent!")

	// disseminate message to all peers
	response := router.disseminate(m)
	if response == nil {
		routerLogger(router, method).Error("Response was nil.")
		return nil
	}

	return response.Params
}

func (router *BroadCastRouter) nextRequestId() int {
	result := router.reqId
	router.reqId++
	return result
}

func (router *BroadCastRouter) disseminate(m *RequestMessage) *ResponseMessage {
	var (
		responseChannel = make(chan *ResponseMessage)
		sendCounter     int
		responseCounter int
	)

	// send message to all registered peerClients---concurrently
	for _, client := range router.routeTable.peers() {
		go router.send(client, m, responseChannel)
		sendCounter += 1
	}

	var result, newMsg *ResponseMessage
	// wait until all client requests responded
	for responseCounter < sendCounter {
		newMsg = <-responseChannel
		responseCounter += 1

		result = router.merge(result, newMsg, m.Method)
	}
	return result
}

func (router *BroadCastRouter) send(client com.RPCClient, request *RequestMessage, responseChannel chan *ResponseMessage) {
	routerLogger(router, request.Method).Debugf("Routing message to client: %s", client.URL())

	// Make the call
	tmpResponse, err := client.Call(router.Name(), *request)
	if tmpResponse != nil && err == nil {
		actualResponse := &ResponseMessage{}
		tmpResponse.Unmarshal(actualResponse)
		responseChannel <- actualResponse
	} else {
		routerLogger(router, request.Method).WithError(err).Errorf("Error while routing message to client: %s", client.URL())
		responseChannel <- nil
	}
}

func (router *BroadCastRouter) merge(oldMsg *ResponseMessage, newMsg *ResponseMessage, handlerName string) *ResponseMessage {

	routerLogger(router, handlerName).Debugf("Merging messages.")

	if oldMsg == nil {
		routerLogger(router, handlerName).Debugf("Old message is nil.")
		return newMsg
	} else if newMsg == nil {
		routerLogger(router, handlerName).Debugf("New message is nil.")
		return oldMsg
	}

	// merge results on the fly; logic for merging has to be implemented by the "upper layer", e.g., the service layer
	if mergeCallback := router.routeHandlers[handlerName].mergeCallback; mergeCallback != nil {
		params := mergeCallback(oldMsg.Params, newMsg.Params)
		return newResponseCopy(params, newMsg)
	}

	return oldMsg
}

/*
Handle a broadcasted message
*/
func (router *BroadCastRouter) Handle(routerName string, msg models.EncapsulatedValues) interface{} {

	var (
		response *ResponseMessage
		request  = RequestMessage{}
	)

	routerLoggerS(router).Debugf("Handle message.", routerName)

	// cast request to RequestMessage from interface and verify it is a valid request
	if err := msg.Unmarshal(&request); err != nil {
		routerLoggerS(router).WithError(err).Infof("Cannot read request while handling message.")
		return nil
	}

	// ignore duplicates to ensure an at least once semantic
	if router.duplicateMap.CheckForDuplicates(request.Src) {
		routerLoggerS(router).WithField("src", request.Src.Id).Info("Dropped duplicate.")
		return nil
	}

	routerLoggerS(router).Infof("No duplicate for: %s", routerName)

	// callback to upper layer
	responseParams := router.deliver(request.Method, request.Params)
	response = newResponse(responseParams, &request)

	// treat every message as BroadcastRouter, therefore:
	// forward message to all other peerClients
	broadcastResponse := router.disseminate(&request)

	// chooseResponse one result
	response = router.merge(response, broadcastResponse, request.Method)

	return *response
}

/*
Route is implemented as 'Broadcast'
*/
func (router *BroadCastRouter) Route(_ Key, method string, message interface{}) (result interface{}) {
	return router.BroadCast(method, message)
}

/*
DelPeer deletes a peer
*/
func (router *BroadCastRouter) DelPeer(key Key) {
	router.routeTable.del(key)
}

/*
NewPeer adds a new peer with the given key and address. This enables the router to forward messages to this peer.
*/
func (router *BroadCastRouter) NewPeer(key Key, address string) {
	if _, found := router.routeTable.get(key); !found {
		routerLoggerS(router).WithField("peer", address).Info("New neighbor.")
		peer := com.NewRPCClient(address)
		router.routeTable.add(key, peer)
	}
}

func (router *BroadCastRouter) deliver(method string, params models.EncapsulatedValues) interface{} {
	if handler, ok := router.routeHandlers[method]; ok {
		routerLogger(router, method).Info("Deliver message to handler.")
		return handler.requestHandler(params)
	} else {
		routerLogger(router, method).Error("Handler not found.")
	}
	return nil
}

/*
Name of the router
*/
func (router *BroadCastRouter) Name() string {
	return router.name
}

/*
AddHandler registers a handler with the router. The handler is called when a message for this handler arrives.
*/
func (router *BroadCastRouter) AddHandler(name string, handler *Handler) {
	routerLoggerS(router).WithField("handler", name).Info("Router registered new callback.")
	router.routeHandlers[name] = handler
}

func newResponseCopy(responseParams interface{}, orig *ResponseMessage) *ResponseMessage {

	if orig == nil {
		return nil
	}

	result, err := NewResponseMessage(orig.RequestSrc, responseParams)
	if err != nil {
		return nil
	}

	return result
}

func newResponse(responseParams interface{}, requestMsg *RequestMessage) (result *ResponseMessage) {

	if responseParams == nil {
		// Error ignored on purpose: result is nil in this error case, which is the expected behaviour of newResponse
		result, _ = NewResponseMessage(requestMsg.Src, "{}")
	} else {
		result, _ = NewResponseMessage(requestMsg.Src, responseParams)
	}

	return result
}
