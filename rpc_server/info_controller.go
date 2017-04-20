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
package rpc_server

//
//import (
//	"net/http"
//	"fmt"
//
//	. "github.com/ottenwbe/golook/app"
//)
//
//const (
//	EP_INFO   = "/info"
//)
//
//
//func init() {
//	HttpServer.RegisterFunction("/", getHome, http.MethodGet)
//	HttpServer.RegisterFunction(EP_INFO, getInfo, http.MethodGet)
//}
//
//
//// Endpoint: GET /
//func getHome(writer http.ResponseWriter, _ *http.Request) {
//	ReturnAck(writer)
//}
//
//// Endpoint: GET /info
//func getInfo(writer http.ResponseWriter, _ *http.Request) {
//	info := NewAppInfo()
//	result := EncodeAppInfo(info)
//	fmt.Fprintln(writer, result)
//}
