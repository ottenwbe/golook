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

	"github.com/ottenwbe/golook/broker/models"

	"github.com/ottenwbe/golook/utils"
	log "github.com/sirupsen/logrus"
	"github.com/ybbus/jsonrpc"
)

type (
	/*JSONRPCClientStub implements the RPCClients interface*/
	JSONRPCClientStub struct {
		serverAddress string
		c             *jsonrpc.RPCClient
	}
	/*JSONRPCReturnType implements the EncapsulatedValues interface*/
	JSONRPCReturnType struct {
		response *jsonrpc.RPCResponse
	}
)

func newJSONRPCClient(address string, port int) RPCClient {

	serverAddress := fmt.Sprintf("http://%s", address)
	if port >= 0 {
		serverAddress = fmt.Sprintf("%s:%d", serverAddress, port)
	}

	return &JSONRPCClientStub{
		serverAddress: serverAddress,
		c:             jsonrpc.NewJsonRPCClient(fmt.Sprintf("%s/rpc", serverAddress)),
	}
}

/*
Call executes a RPC call with the given method name and the given parameters.
It returns a generic return value. Clients can retrieve the result by calling the Unmarshal method on the result.
*/
func (lc *JSONRPCClientStub) Call(method string, parameters interface{}) (models.EncapsulatedValues, error) {

	log.WithField("method", method).WithField("com", "jsonrpc").Debugf("Making a call for %s with %s", method, utils.MarshalSD(parameters))

	response, err := lc.c.Call(method, parameters)
	if err != nil {
		return nil, err
	}

	r := &JSONRPCReturnType{response: response}
	log.WithField("method", method).WithField("com", "jsonrpc").Debugf("Getting a response for %s with %s", method, r.response)

	return r, nil
}

/*
URL returns the url of the RPC server to which this client connects
*/
func (lc *JSONRPCClientStub) URL() string {
	return lc.serverAddress
}

/*
Unmarshal allows callers of the RPC client to unmarshal the result retrieved by the client
*/
func (rt *JSONRPCReturnType) Unmarshal(v interface{}) error {
	return rt.response.GetObject(v)
}
