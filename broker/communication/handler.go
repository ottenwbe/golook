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
	"errors"
	"github.com/ottenwbe/golook/broker/models"
	log "github.com/sirupsen/logrus"
)

type (
	/*MessageHandler is the interface that needs to be implemented to receive messages from the communication layer.*/
	MessageHandler interface {
		Handle(method string, params models.EncapsulatedValues) interface{}
	}

	dispatcherBinding struct {
		handler  MessageHandler
		receiver RPCServer
	}

	dispatcherBindings map[string]dispatcherBinding
)

/*
MessageDispatcher allows components to register handlers that are called when messages arrive for them.
In turn, server components (e.g., the 'json_rpc_server') can delegate the job of dispatching messages to MessageDispatcher.
*/
var MessageDispatcher = newMessageDispatcher()

func (dispatcherBindings *dispatcherBindings) handleMessage(router string, message models.EncapsulatedValues) (interface{}, error) {
	if reg, ok := (*dispatcherBindings)[router]; ok && reg.handler != nil {
		return reg.handler.Handle(router, message), nil
	}
	log.Info("Method dropped before handing it over to handler. No handler registered.")
	return nil, errors.New("Method dropped before handing it over to handler. No handler registered")
}

/*
RegisterHandler takes (message) handler with a given name, (expected) request type, and (expected) response type as input.
The registered handlers are called when messages for a given 'name' arrive.
Note: When a handler is registered you have to remove the handler again with 'RemoveHandler' to ensure that no messages is dispatched to the handler.
*/
func (dispatcherBindings *dispatcherBindings) RegisterHandler(name string, handler MessageHandler, requestType interface{}, responseType interface{}) {
	receiver := newRPCServer(name)
	(*dispatcherBindings)[name] = dispatcherBinding{handler, receiver}
	receiver.Associate(name, requestType, responseType)
}

/*
RemoveHandler ensures that a given handler no longer receives messages.
*/
func (dispatcherBindings *dispatcherBindings) RemoveHandler(name string) {
	if e, ok := (*dispatcherBindings)[name]; ok {
		delete(*(dispatcherBindings), name)
		e.receiver.Finalize()
	}
}

/*
HasHandler checks if a handler with a given 'name' is managed by the collection of dispatcherBindings 'r'
*/
func (dispatcherBindings *dispatcherBindings) HasHandler(name string) bool {
	_, ok := (*dispatcherBindings)[name]
	return ok
}

func newMessageDispatcher() *dispatcherBindings {
	tmp := make(dispatcherBindings)
	return &tmp
}
