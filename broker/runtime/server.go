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

	"errors"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"net"
	"reflect"
	"sync"
)

type (
	Server interface {
		StartServer(wg *sync.WaitGroup)
		Stop()
		Info() map[string]interface{}
		Name() string
		IsRunning() bool
	}
	HTTPSever struct {
		Address  string
		server   *http.Server
		router   http.Handler
		listener net.Listener
	}
)

type ServerType string

const (
	ServerHttp ServerType = "http"
)

var (
	servers = []Server{}
)

func NewServer(address string, serverType ServerType) (server Server) {
	switch serverType {
	case ServerHttp:
		server = &HTTPSever{
			server:  nil,
			router:  mux.NewRouter().StrictSlash(true),
			Address: address,
		}
	default:
		log.WithField("type", string(serverType)).Error("Server could not be created. Type is not supported.")
	}

	if server != nil {
		servers = append(servers, server)
	}

	return server
}

func (s *HTTPSever) Name() string {
	return s.Address
}

func (s *HTTPSever) IsRunning() bool {
	return s.server != nil
}

func (s *HTTPSever) StartServer(wg *sync.WaitGroup) {
	if wg == nil || s.router == nil || s.Address == "" {
		log.Error("Http server cannot be started. Please ensure that all parameters are set: Address and router. Moreover, ensure that a wg is provided during startup.")
	}

	defer wg.Done()

	s.server = &http.Server{Addr: s.Address, Handler: s.router}

	var err error
	// make our own listener
	s.listener, err = net.Listen("tcp", s.Address)
	if err != nil {
		log.Panic(err)
	}

	// start the httpServer and listen
	log.Error(s.server.Serve(s.listener))
}

func (s *HTTPSever) Stop() {
	s.listener.Close()
	s.server = nil
}

/*
RegisterFunction registers an endpoint and the corresponding controller for a specific http method.
*/
func (s *HTTPSever) RegisterFunction(path string, f func(http.ResponseWriter, *http.Request), methods ...string) error {
	log.Infof("Register http endpoint: %s", path)
	if router, ok := s.router.(*mux.Router); ok {
		router.HandleFunc(path, f).Methods(methods...)
		return nil
	}
	return errors.New("Cannot register http function")

}

func (s *HTTPSever) Info() map[string]interface{} {
	info := map[string]interface{}{}

	info["endpoints"] = s.RegisteredEndpoints()
	serverInfo := map[string]interface{}{}
	if s.server != nil {
		serverInfo["address"] = s.server.Addr
		serverInfo["readTimeOut"] = s.server.ReadTimeout
		serverInfo["writeTimeOut"] = s.server.WriteTimeout
		serverInfo["maxHeaderBytes"] = s.server.MaxHeaderBytes
	} else {
		serverInfo["address"] = s.Address
	}
	serverInfo["router"] = reflect.TypeOf(s.router).String()
	info["server"] = serverInfo

	return info
}

/*
RegisteredEndpoints returns all registered endpoints as an array of string.

Example result:
["/info","/system/{system}","/foo/bar"]
*/
func (s *HTTPSever) RegisteredEndpoints() []string {
	result := []string{}

	if router, ok := s.router.(*mux.Router); ok {
		router.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
			if s, err := route.GetPathTemplate(); err == nil {
				result = append(result, s)
			}
			return nil
		})
	}

	return result
}

/**
RunServer starts all server and waits for them to return, i.e. they gracefully shut down.
*/
func RunServer() {
	var wg = &sync.WaitGroup{}

	wg.Add(len(servers))

	for _, server := range servers {
		go server.StartServer(wg)
	}

	wg.Wait()
}

func StopServer() {
	for _, server := range servers {
		go server.Stop()
	}
}

/**
GetServer returns a list of registered servers
*/
func ServerState() map[string]bool {
	result := map[string]bool{}
	for _, server := range servers {
		result[reflect.TypeOf(server).String()+"/"+server.Name()] = server.IsRunning()
	}
	return result
}
