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
	"bytes"
	"github.com/gorilla/mux"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

//Verify that home exists and returns the correct status code
func TestHome(t *testing.T) {

	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(home)

	handler.ServeHTTP(rr, req)

	// Expect status 200
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	expected := `{up}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestSystemLifecycle(t *testing.T) {
	postTestSystem(t)
	getTestSystem(t)
	delTestSystem(t)
}

func postTestSystem(t *testing.T) {
	var jsonStr = []byte(`{"name":"1234","os":"linux","ip":"1.1.1.1","uuid":"a"}`)
	req, err := http.NewRequest("POST", "/systems/linux-test", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(PostSystem)

	handler.ServeHTTP(rr, req)

	// Expect status 200
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	expected := "{\"id\":"
	if !strings.Contains(rr.Body.String(), expected) {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}

}

func getTestSystem(t *testing.T) {
	req, err := http.NewRequest("GET", "/systems/a", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	m := mux.NewRouter()
	m.HandleFunc("/systems/{system}", GetSystem)
	m.ServeHTTP(rr, req)

	// Check the response body is what we expect.
	expected := "\"ip\":\"1.1.1.1\""
	if !strings.Contains(rr.Body.String(), expected) {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}

}

func delTestSystem(t *testing.T) {
	req, err := http.NewRequest("DELETE", "/systems/a", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	m := mux.NewRouter()
	m.HandleFunc("/systems/{system}", DelSystem)
	m.ServeHTTP(rr, req)

	// Check the response body is what we expect.
	expected := "Deleting"
	if !strings.Contains(rr.Body.String(), expected) {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}