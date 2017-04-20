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
	"github.com/ottenwbe/golook/broker/models"
	"github.com/ottenwbe/golook/broker/runtime"
)

type (
	RpcHandler struct{}
	RpcResult  struct {
		Message string `json:"message"`
	}
)

var _ (jsonrpc.Handler) = (*RpcHandler)(nil)

func (h *RpcHandler) ServeJSONRPC(c context.Context, params *json.RawMessage) (interface{}, *jsonrpc.Error) {

	var p models.RequestMessage
	if err := jsonrpc.Unmarshal(params, &p); err != nil {
		return nil, err
	}

	return ToRouteLayer(p.Method, p.Content), nil

}

func init() {
	runtime.HttpServer.RegisterFunctionS("/rpc", jsonrpc.HandlerFunc)
	runtime.HttpServer.RegisterFunctionS("/rpc/debug", jsonrpc.DebugHandlerFunc)

	jsonrpc.RegisterMethod("encapsulated", &RpcHandler{}, models.RequestMessage{}, models.ResponseMessage{})
}
