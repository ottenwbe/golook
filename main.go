//Copyright 2016 Beate Ottenwälder
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
	"time"
)

//todo: refactor and move to own file
type File struct {
	Name    string    `json:"name"`
	Created time.Time `json:"created"`
}

type System struct {
	Name  string `json:"name"`
	OS    string `json:"os"`
	IP    string `json:"ip"`
	Files []File `json:"files"`
}

type repository interface {
	findFile(fileName string) ([]System, error)
	get(fileName string, systemName string) (File, error)
	putSystem(system System) error
	putFile(file File) error
}

//var s = make([]System, 1)

func main() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", home)
	router.HandleFunc("/files/{file}", getFile).Methods("GET")
	router.HandleFunc("/files/{file}", putFile).Methods("PUT")
	router.HandleFunc("/systems/{system}", getSystem).Methods("GET")
	router.HandleFunc("/systems/{system}", putSystem).Methods("PUT")
	log.Fatal(http.ListenAndServe(":8080", router))
}

func home(writer http.ResponseWriter, request *http.Request) {
	fmt.Fprintf(writer, "Remote File Lookup")
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

func putSystem(writer http.ResponseWriter, request *http.Request) {
	var system System
	err := json.NewDecoder(request.Body).Decode(&system)
	//append(s, system)
	if err != nil {
		fmt.Fprintln(writer, "Error during conversion")
	} else {
		fmt.Fprintln(writer, "List of files that are found for:", system.Name)
	}
}
