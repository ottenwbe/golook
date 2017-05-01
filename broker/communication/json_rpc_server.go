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
	"context"
	"encoding/json"

	"errors"
	"github.com/osamingo/jsonrpc"
	"github.com/ottenwbe/golook/broker/runtime"
	log "github.com/sirupsen/logrus"
)

type (
	JsonRPCServerStub struct {
		handler string
		active  bool
	}

	JsonRPCParams struct {
		params json.RawMessage
	}
)

var (
	HttpRpcServer runtime.Server
)

var _ (jsonrpc.Handler) = (*JsonRPCServerStub)(nil)

func (rpc *JsonRPCServerStub) ServeJSONRPC(_ context.Context, params *json.RawMessage) (interface{}, *jsonrpc.Error) {

	// if the interface is not active, return an error
	if !rpc.active {
		return nil, jsonrpc.ErrMethodNotFound()
	}

	log.Info("Received RPC message: %s", string(*params))

	p := &JsonRPCParams{params: *params}

	if response, err := MessageDispatcher.handleMessage(rpc.handler, p); err != nil {
		log.WithError(err).Error("Error when dispatching Json RPC call")
		return response, jsonrpc.ErrMethodNotFound()
	} else {
		return response, nil
	}
}

func (rpc *JsonRPCServerStub) Associate(handlerName string, request interface{}, response interface{}) {
	rpc.handler = handlerName
	rpc.active = true
	jsonrpc.RegisterMethod(handlerName, rpc, request, response)
}

func (rpc *JsonRPCServerStub) Finalize() {
	//Unfortuanetely jsonrpc has no method for removing a registered function, therefore we only mark it as deleted
	rpc.active = false
}

func (p *JsonRPCParams) Unmarshal(v interface{}) error {

	var interfaceParams []json.RawMessage
	if err := jsonrpc.Unmarshal(&p.params, &interfaceParams); err != nil {
		return err
	}

	if len(interfaceParams) == 1 {
		if err := jsonrpc.Unmarshal(&interfaceParams[0], v); err != nil {
			return err
		}
	} else {
		return errors.New("Slices are not supported")

	}

	return nil
}
