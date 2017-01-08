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
	"time"
	"io"
	"crypto/rand"
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

var s = make(map[string]System, 1)

func main() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", home)
	router.HandleFunc("/files/{file}", getFile).Methods("GET")
	router.HandleFunc("/files/{file}", putFile).Methods("PUT")
	router.HandleFunc("/systems/{system}", getSystem).Methods("GET")
	router.HandleFunc("/systems", postSystem).Methods("POST")
	//router.HandleFunc("/systems/{system}", postSystem).Methods("PUT")
	router.HandleFunc("/systems/{system}", delSystem).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":8080", router))
}

func home(writer http.ResponseWriter, request *http.Request) {
	fmt.Fprintf(writer, "Remote File Lookup is up")
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

//Source: https://play.golang.org/p/4FkNSiUDMg
func newUUID() (string, error) {
	uuid := make([]byte, 16)
	n, err := io.ReadFull(rand.Reader, uuid)
	if n != len(uuid) || err != nil {
		return "", err
	}
	// variant bits; see section 4.1.1
	uuid[8] = uuid[8] &^ 0xc0 | 0x80
	// version 4 (pseudo-random); see section 4.1.3
	uuid[6] = uuid[6] &^ 0xf0 | 0x40
	return fmt.Sprintf("%x-%x-%x-%x-%x", uuid[0:4], uuid[4:6], uuid[6:8], uuid[8:10], uuid[10:]), nil
}

func postSystem(writer http.ResponseWriter, request *http.Request) {
	var system System
	err := json.NewDecoder(request.Body).Decode(&system)
	newSystem, _ := newUUID()
	s[newSystem] = system
	if err != nil {
		fmt.Fprintln(writer, "Error during conversion")
	} else {
		fmt.Fprintln(writer, "List of files that are found for:", system.Name)
	}
}
