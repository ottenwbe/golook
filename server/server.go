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
package server

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

var httpServer *http.Server

func createRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", home).Methods("GET").Name("/")
	router.HandleFunc("/files/{file}", getFile).Methods("GET")
	router.HandleFunc("/systems/{system}/files/{file}", getSystemFile).Methods("GET")
	router.HandleFunc("/systems/{system}/files/{file}", putFile).Methods("PUT")
	router.HandleFunc("/systems/{system}/files/{file}", putFiles).Methods("PUT")
	router.HandleFunc("/systems/{system}", getSystem).Methods("GET")
	router.HandleFunc("/systems", putSystem).Methods("PUT")
	router.HandleFunc("/systems/{system}", delSystem).Methods("DELETE")
	return router
}

func StartServer(address string) {
	router := createRouter()
	httpServer = &http.Server{Addr: address, Handler: router}
	// start the httpServer and listen
	log.Fatal(httpServer.ListenAndServe())
}

//func StopServer() error {
//TODO: wait for graceful shutdown in go 1.8
//}
