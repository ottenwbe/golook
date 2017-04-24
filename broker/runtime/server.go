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
package runtime

import (
	"net/http"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

type lookSever struct {
	Address string
	server  *http.Server
	router  *mux.Router
}

var (
	HttpServer *lookSever = &lookSever{
		server:  nil,
		router:  mux.NewRouter().StrictSlash(true),
		Address: "",
	}
)

func (s *lookSever) StartServer() {
	s.server = &http.Server{Addr: s.Address, Handler: s.router}
	// start the httpServer and listen
	log.Fatal(s.server.ListenAndServe())
}

//func (s *lookSever) StopServer() error {
//TODO: wait for graceful shutdown in go 1.8
// 	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
//	defer cancel()
//	s.server.Shutdown(ctx)
//}

/*
	Register an endpoint and the corresponding controller.
*/
func (s *lookSever) RegisterFunctionS(path string, f func(http.ResponseWriter, *http.Request)) {
	log.Infof("Register http endpoint: %s", path)
	http.HandleFunc(path, f)
}

/*
	Register an endpoint and the corresponding controller for a specific http method.
*/
func (s *lookSever) RegisterFunction(path string, f func(http.ResponseWriter, *http.Request), method string) {
	log.Infof("Register http endpoint: %s", path)
	s.router.HandleFunc(path, f).Methods(method)
}

/*
	Returns all registered endpoints as an array of string.

	Example result:
	["/info","/system/{system}","/foo/bar"]
*/
func (s *lookSever) RegisteredEndpoints() []string {
	result := []string{}

	s.router.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		if s, err := route.GetPathTemplate(); err == nil {
			result = append(result, s)
		}
		return nil
	})

	return result
}
