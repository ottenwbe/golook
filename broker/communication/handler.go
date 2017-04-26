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
)

type (
	MessageHandler interface {
		Handle(method string, params models.EncapsulatedValues) interface{}
	}

	DispatcherBinding struct {
		handler  MessageHandler
		receiver RpcServer
	}

	DispatcherBindings map[string]DispatcherBinding
)

func (r *DispatcherBindings) handleMessage(router string, message models.EncapsulatedValues) (interface{}, error) {
	if reg, ok := (*r)[router]; ok && reg.handler != nil {
		return reg.handler.Handle(router, message), nil
	} else {
		return nil, errors.New("Method dropped before handing it over to handler. No handler registered.")
	}

}

func (r *DispatcherBindings) RegisterHandler(handlerName string, handler MessageHandler, requestType interface{}, responseType interface{}) {
	receiver := newRPCServer(handlerName)
	(*r)[handlerName] = DispatcherBinding{handler, receiver}
	receiver.Associate(handlerName, requestType, responseType)
}

func (r *DispatcherBindings) RemoveHandler(name string) {
	if e, ok := (*r)[name]; ok {
		delete(*(r), name)
		e.receiver.Finalize()
	}
}

func (r *DispatcherBindings) HasHandler(name string) bool {
	_, ok := (*r)[name]
	return ok
}

var MessageDispatcher = newMessageDispatcher()

func newMessageDispatcher() *DispatcherBindings {
	tmp := make(DispatcherBindings)
	return &tmp
}
