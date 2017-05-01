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
	"net/http"

	jsonrpcServer "github.com/osamingo/jsonrpc"
	"github.com/spf13/viper"

	golook "github.com/ottenwbe/golook/broker/runtime"
)

func ApplyCommunicationConfiguration() {

	ClientType = viper.GetString("communication.client.type")
	port = viper.GetInt("communication.jsonrpc.client.port")

	serverType = viper.GetString("communication.server.type")
	HttpRpcServer = golook.GetOrCreate(viper.GetString("communication.jsonrpc.server.address"), golook.ServerHttp)
	HttpRpcServer.(*golook.HTTPSever).RegisterFunction("/rpc", jsonrpcServer.HandlerFunc, http.MethodGet, http.MethodPost)
	HttpRpcServer.(*golook.HTTPSever).RegisterFunction("/rpc/debug", jsonrpcServer.DebugHandlerFunc, http.MethodGet, http.MethodPost)
}

func InitCommunicationConfiguration() {

	viper.SetDefault("communication.client.type", jsonRPC)
	viper.SetDefault("communication.server.type", jsonRPC)
	viper.SetDefault("communication.jsonrpc.client.port", 8382)
	viper.SetDefault("communication.jsonrpc.server.address", ":8382")
}
