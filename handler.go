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
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func home(writer http.ResponseWriter, request *http.Request) {
	fmt.Fprint(writer, "{up}")
}

//todo: return list of files that match the search pattern
func getFile(writer http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	fmt.Fprintln(writer, "List of files that are found for:", params["file"])
}

func putFile(writer http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	fmt.Fprintln(writer, "List of files that are found for:", params["file"])
}

//todo: return list of files that match the search pattern
func getSystem(writer http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	fmt.Fprintln(writer, "List of files that are found for:", params["system"])
}

func delSystem(writer http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	delete(s, params["system"])
	fmt.Fprintln(writer, "Deleting system:", params["system"])
}

func postSystem(writer http.ResponseWriter, request *http.Request) {
	if (request != nil) && (request.Body != nil) {
		var system System
		err := json.NewDecoder(request.Body).Decode(&system)
		if err != nil {
			log.Print("Error: Post system request was error")
			fmt.Fprint(writer, "Post system request was error")
		} else {
			extractSystem(&system, &writer)
		}
	} else {
		log.Print("Error: Post system request was empty")
		fmt.Fprint(writer, "Post system request was empty")
	}
}

func extractSystem(system *System, writer *http.ResponseWriter) {
	newSystem, erruuid := NewUUID()
	if erruuid == nil {
		s[newSystem] = *system
		log.Printf("Post system request was a success %s", *system)
		fmt.Fprint(*writer, newSystem)
	} else {
		log.Printf("Error: UUID error %s", erruuid)
		fmt.Fprint(*writer, "UUID error")
	}
}
