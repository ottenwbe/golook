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
	address string
	server *http.Server
	router *mux.Router
}

var (
	HttpServer *lookSever
)

func (s *lookSever) StartServer() {
	s.server = &http.Server{Addr: s.address, Handler: s.router}
	// start the httpServer and listen
	log.Fatal(s.server.ListenAndServe())
}

//func (s *lookSever) StopServer() error {
//TODO: wait for graceful shutdown in go 1.8
// 	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
//	defer cancel()
//	s.server.Shutdown(ctx)
//}

func (s *lookSever) RegisterFunction(path string, f func(http.ResponseWriter, *http.Request), method string) {
	s.router.HandleFunc(path, f).Methods(method)
}

func init() {
	HttpServer = &lookSever{
		server: nil,
		router: mux.NewRouter().StrictSlash(true),
		address: "",
	}

	//NOTE: Requires that HTTPServer is instantiated!
	RootCmd.Flags().StringVarP(&HttpServer.address, "httpserver", "s", ":8383", "(optional) Default address of the http server. Default: ':8383'")
}

