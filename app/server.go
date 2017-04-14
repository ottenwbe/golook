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
package app

import (
	"net/http"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

type lookSever struct {
	server *http.Server
	router *mux.Router
}

var (
	HttpServer *lookSever
)

func (s *lookSever) StartServer(address string) {
	s.server = &http.Server{Addr: address, Handler: s.router}
	// start the httpServer and listen
	log.Fatal(s.server.ListenAndServe())
}

//func (s *lookSever) StopServer() error {
//TODO: wait for graceful shutdown in go 1.8
//}

func (s *lookSever) RegisterFunction(path string, f func(http.ResponseWriter, *http.Request), method string) {
	s.router.HandleFunc(path, f).Methods(method)
}

func init() {
	HttpServer = &lookSever{
		server: nil,
		router: mux.NewRouter().StrictSlash(true),
	}
}
