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
	"github.com/ybbus/jsonrpc"

	. "github.com/ottenwbe/golook/broker/models"

	"fmt"
)

// implements the LookupClient interface
type LookupRPCClient struct {
	serverUrl string
	c         *jsonrpc.RPCClient //TODO: check if http communication is synchronized
}

func init() {
	LookupClients.Add("rpc", NewLookupRPCClient, true)
}

func NewLookupRPCClient(url string) LookupClient {
	return &LookupRPCClient{
		serverUrl: url,
		c:         jsonrpc.NewRPCClient(fmt.Sprintf("%s/rpc", url)),
	}
}

func (lc *LookupRPCClient) Call(method string, message interface{}) (interface{}, error) {
	m, err := NewRpcMessage(method, message)
	if err != nil {
		return nil, err
	}

	response, err := lc.c.Call("encapsulated", m)
	if err != nil {
		return nil, err
	}

	var msg ResponseMessage
	err = response.GetObject(&msg)
	if err != nil {
		return nil, err
	}

	return msg.Content, nil
}
