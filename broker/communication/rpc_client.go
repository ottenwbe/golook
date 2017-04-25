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
	"fmt"

	. "github.com/ottenwbe/golook/broker/models"

	"github.com/sirupsen/logrus"
	"github.com/ybbus/jsonrpc"
)

// implements the LookupClient interface
type LookupRPCClient struct {
	serverUrl string
	c         *jsonrpc.RPCClient
}

func NewLookupRPCClient(url string) LookupClient {
	return &LookupRPCClient{
		serverUrl: url,
		c:         jsonrpc.NewRPCClient(fmt.Sprintf("%s/rpc", url)),
	}
}

type RPCMsgResponse struct {
	response *jsonrpc.RPCResponse
}

func (r *RPCMsgResponse) GetObject(v interface{}) error {
	return r.response.GetObject(v)
}

//func (lc *LookupRPCClient) Call(index string, method string, message interface{}) (interface{}, error) {
func (lc *LookupRPCClient) Call(method string, m interface{}) (MsgParams, error) {

	logrus.Infof("Making a call for %s with %s", method, m)

	response, err := lc.c.Call(method, m)
	if err != nil {
		return nil, err
	}

	r := &RPCMsgResponse{response: response}

	logrus.Info("response is " + string(response.ID))

	return r, nil
}

func (lc *LookupRPCClient) Url() string {
	return lc.serverUrl
}
