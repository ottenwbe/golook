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

const (
	nack = "{nack}"
	ack  = "{ack}"

	systemPath = "system"
	filePath   = "file"
)

/////////////////////////////////////
// Controllers
/////////////////////////////////////

func home(writer http.ResponseWriter, _ *http.Request) {
	fmt.Fprint(writer, ack)
}

func getFile(writer http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)

	var result map[string]*System

	findString := params[filePath]

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

	//TODO error handling
	bytes, _ := json.Marshal(systemMap[params[systemPath]].Files)
	fmt.Fprintln(writer, string(bytes))
}

func getSystemFile(writer http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)

	if sys, ok := systemMap[params[systemPath]]; ok {
		bytes, marshallErr := json.Marshal(sys.Files)
		if marshallErr == nil {
			fmt.Fprintln(writer, string(bytes))
		} else {
			fmt.Fprintln(writer, nack)
			log.Printf("Error marshalling file %s", marshallErr)
		}
	} else {
		fmt.Fprintln(writer, nack)
		log.Print("System to receive file for cannot be found")
	}
}

func putFiles(writer http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)

	var files []File
	_ = json.NewDecoder(request.Body).Decode(&files)
	systemMap[params[systemPath]].Files = files

	fmt.Fprintln(writer, "Accepted")
}

func putFile(writer http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request) //params

	var file File
	err := json.NewDecoder(request.Body).Decode(&file)
	if err != nil {
		fmt.Fprintln(writer, nack)
		log.Printf("Error. File could not be decoded while putting the file to server %s", err)
	} else if sys, ok := systemMap[params[systemPath]]; ok {
		sys.Files = append(sys.Files, file)
		fmt.Fprintln(writer, "{ack}")
	} else {
		fmt.Fprintln(writer, nack)
		log.Printf("Error. System (%s) not found while putting a file information to the server.", params[systemPath])
	}
}

func getSystem(writer http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)

	if sys, ok := systemMap[params[systemPath]]; ok {
		log.Printf("Find (system=%s) in '%d' systems", params[systemPath], len(systemMap))
		str, marshalError := json.Marshal(sys)
		if marshalError != nil {
			log.Print("Json could not be marshalled")
			fmt.Fprint(writer, "{}")
		} else {
			log.Printf("Find (system=%s) was successful: %s", params[systemPath], str)
			fmt.Fprint(writer, string(str))
		}
	} else {
		log.Printf("System (%s) not found: Returning empty Json", params[systemPath])
		fmt.Fprint(writer, "{}")
	}
}

func delSystem(writer http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	delete(systemMap, params[systemPath])
	fmt.Fprintln(writer, "Deleting system:", params[systemPath])
}

func putSystem(writer http.ResponseWriter, request *http.Request) {
	if (request != nil) && (request.Body != nil) {
		tryAddSystem(&writer, request)
	} else {
		log.Print("Error: Post system request was empty")
		fmt.Fprint(writer, nack)
	}
}

/////////////////////////////////////
// Helpers for controllers
/////////////////////////////////////

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
		log.Printf("Post (system=%s) request was a success %s", newSystem, systemMap[newSystem])
		fmt.Fprint(*writer, "{\"id\":\""+newSystem+"\"}")
	} else {
		log.Printf("Error: UUID error %s", errUuid)
		fmt.Fprint(*writer, "UUID error")
	}
}
