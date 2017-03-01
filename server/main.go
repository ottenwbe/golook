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
package main

import (
	"log"

	"github.com/gorilla/mux"
	"net/http"
)

func makeServer() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", home)
	router.HandleFunc("/files/{file}", GetFile).Methods("GET")
	router.HandleFunc("/systems/{system}/files/{file}", GetSystemFile).Methods("GET")
	router.HandleFunc("/systems/{system}/files/{file}", PutFile).Methods("PUT")
	router.HandleFunc("/systems/{system}/files/{file}", PutFiles).Methods("PUT")
	router.HandleFunc("/systems/{system}", GetSystem).Methods("GET")
	router.HandleFunc("/systems", PostSystem).Methods("POST")
	router.HandleFunc("/systems/{system}", DelSystem).Methods("DELETE")

	// start the server and listen on port 8080
	log.Fatal(http.ListenAndServe(":8080", router))
}

func main() {
	makeServer()
}
