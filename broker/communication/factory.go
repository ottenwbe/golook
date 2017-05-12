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

const (
	jsonRPC = "jsonrpc"
	/*MockRPC represents the name of the configuration where mocked RPC clients are created*/
	MockRPC = "mockrpc"
)

var (
	// value is injected through configuration (see configuration.go)
	serverType string
	/*ClientType informs NewRPCClient what type of client should be built.
	Note: the actual value is injected through configuration (see configuration.go)*/
	ClientType string

	port = 8382
)

func newRPCServer(associatedHandler string) RPCServer {
	return &JSONRPCServerStub{handler: associatedHandler}
}

/*
NewRPCClient returns a new RPCClient. The actual type of the RPCClient needs to be configured beforehand, e.g., in a config file.
*/
func NewRPCClient(url string) RPCClient {
	switch ClientType {
	case MockRPC:
		return newMockClient()
	default:
		return newJSONRPCClient(url, port)
	}

}
