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
	"github.com/ottenwbe/golook/helper"
	"log"
	"net/http"
	"strings"
)

func home(writer http.ResponseWriter, _ *http.Request) {
	fmt.Fprint(writer, "{up}")
}

func GetFile(writer http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)

	var result map[string]*System

	findString := params["file"]

	for sid, system := range systemMap {
		for _, file := range system.Files {
			if strings.Contains(file.Name, findString) {
				if _, ok := result[sid]; !ok {
					result[sid] = new(System)
					result[sid].Name = system.Name
				}
				result[sid].Files = append(result[sid].Files, file)
			}
		}
	}

	bytes, _ := json.Marshal(systemMap[params["system"]].Files)
	fmt.Fprintln(writer, string(bytes))
}

func GetSystemFile(writer http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)

	bytes, _ := json.Marshal(systemMap[params["system"]].Files)
	fmt.Fprintln(writer, string(bytes))
}

func PutFiles(writer http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)

	var files []File
	_ = json.NewDecoder(request.Body).Decode(&files)
	systemMap[params["system"]].Files = files

	fmt.Fprintln(writer, "Accepted")
}

func PutFile(writer http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)

	var file File
	_ = json.NewDecoder(request.Body).Decode(&file)
	systemMap[params["system"]].Files = append(systemMap[params["system"]].Files, file)

	fmt.Fprintln(writer, "Accepted")
}

func GetSystem(writer http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	fmt.Fprintln(writer, "List of files that are found for:", params["system"])
	log.Printf("Get (system=%s) from '%d' systems", params["system"], len(systemMap))

	if sys, ok := params["system"]; ok {
		str, _ := json.Marshal(systemMap[sys]) //TODO: error handling
		fmt.Fprint(writer, string(str))
	} else {
		log.Print("Error: Get system failed")
		fmt.Fprint(writer, "System not found")
	}
}

func DelSystem(writer http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	delete(systemMap, params["system"])
	fmt.Fprintln(writer, "Deleting system:", params["system"])
}

func PostSystem(writer http.ResponseWriter, request *http.Request) {
	if (request != nil) && (request.Body != nil) {
		tryAddSystem(&writer, request)
	} else {
		log.Print("Error: Post system request was empty")
		fmt.Fprint(writer, "Post system request was empty")
	}
}

func tryAddSystem(writer *http.ResponseWriter, request *http.Request) {
	var system System
	err := json.NewDecoder(request.Body).Decode(&system)
	if err != nil {
		log.Printf("Error: Post system request has errors: %s", err)
		fmt.Fprint(*writer, "Post system request has errors")
	} else {
		extractSystem(&system, writer)
	}
}

func extractSystem(system *System, writer *http.ResponseWriter) {
	var newSystem string
	var errUuid error = nil
	if system.UUID == "" {
		newSystem, errUuid = helper.NewUUID()
	} else {
		newSystem = system.UUID
	}
	if errUuid == nil {
		systemMap[newSystem] = system
		log.Printf("Post system request was a success %s", systemMap[newSystem])
		fmt.Fprint(*writer, "{\"id\":\"" + newSystem + "\"}")
	} else {
		log.Printf("Error: UUID error %s", errUuid)
		fmt.Fprint(*writer, "UUID error")
	}
}
