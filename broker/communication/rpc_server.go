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
	"github.com/osamingo/jsonrpc"
	"github.com/ottenwbe/golook/broker/runtime"
)

type (
	RpcHandler struct {
		index string
	}

	RpcReceiverParams struct {
		params *json.RawMessage
	}
)

func (p *RpcReceiverParams) GetObject(v interface{}) error {
	if err := jsonrpc.Unmarshal(p.params, v); err != nil {
		return err
	}
	return nil
}

var _ (jsonrpc.Handler) = (*RpcHandler)(nil)

func (h *RpcHandler) ServeJSONRPC(_ context.Context, params *json.RawMessage) (interface{}, *jsonrpc.Error) {

	var p = &RpcReceiverParams{params: params}

	return toRouteLayer(h.index, p), nil
}

func init() {
	runtime.HttpServer.RegisterFunctionS("/rpc", jsonrpc.HandlerFunc)
	runtime.HttpServer.RegisterFunctionS("/rpc/debug", jsonrpc.DebugHandlerFunc)
}
