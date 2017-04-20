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
package api

/*
	Common (helper) functions and constants, required by all controllers.
*/

import (
	"errors"
	"fmt"
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/gorilla/mux"
)

const (
	Nack = "{Nack}"
	Ack  = "{Ack}"

	systemPath = "system"
	filePath   = "file"
)

/*
	Helpers Functions for controller
*/

func IsValidRequest(request *http.Request) bool {
	return (request != nil) && (request.Body != nil)
}

func ReturnAck(writer http.ResponseWriter) (int, error) {
	return fmt.Fprint(writer, Ack)
}

func ReturnNackAndLog(writer http.ResponseWriter, errorString string, status int) {
	log.Print(errorString)
	http.Error(writer, errors.New(Nack).Error(), status)
}

func ReturnNackAndLogError(writer http.ResponseWriter, errorString string, err error, status int) {
	log.WithError(err).Print(errorString)
	http.Error(writer, errors.New(Nack).Error(), status)
}

//func extractSystemFromPath(request *http.Request) string {
//	params := mux.Vars(request)
//	system := params[systemPath]
//	return system
//}

func extractFileFromPath(request *http.Request) string {
	params := mux.Vars(request)
	fileName := params[filePath]
	return fileName
}
