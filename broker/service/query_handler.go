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

package service

import (
	"github.com/ottenwbe/golook/broker/models"
	. "github.com/ottenwbe/golook/broker/repository"
	. "github.com/ottenwbe/golook/broker/utils"
	log "github.com/sirupsen/logrus"
)

const (
	FILE_QUERY = "file query"
)

func handleFileQuery(params models.EncapsulatedValues) interface{} {

	var (
		systemMessage PeerFileQuery
		response      PeerResponse
	)

	err := params.Unmarshal(&systemMessage)

	if err == nil {
		response = processFileQuery(systemMessage)
	} else {
		log.WithError(err).Error("Could not handle file query")
		response = PeerResponse{Error: true, Message: "Problem Unmarshalling file query", Data: nil}
	}

	return response
}

func processFileQuery(systemMessage PeerFileQuery) PeerResponse {
	result := GoLookRepository.FindSystemAndFiles(systemMessage.SearchString)
	tmp, err := MarshalB(result)
	return PeerResponse{Error: err != nil, Message: "Processed file query", Data: tmp}
}
