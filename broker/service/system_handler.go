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
	"fmt"
	. "github.com/ottenwbe/golook/broker/repository"
	. "github.com/ottenwbe/golook/broker/utils"
	log "github.com/sirupsen/logrus"
)

const (
	SYSTEM_REPORT = "system_report"
)

func handleSystemReport(params interface{}) interface{} {
	log.Info("handleSystemReport !!!")
	var (
		systemMessage SystemTransfer
		response      PeerResponse
	)

	log.Info("handleSystemReport will gather for %s", params.(string))

	if err := Unmarshal(params, &systemMessage); err == nil {
		if systemMessage.IsDeletion {
			GoLookRepository.DelSystem(systemMessage.Uuid)
			response = PeerResponse{false, fmt.Sprintf("Processed request for deleting system %s", systemMessage.Uuid), nil}
		} else {
			response.Error = !GoLookRepository.StoreSystem(systemMessage.Uuid, systemMessage.System)
			response.Message = fmt.Sprintf("Processed request for adding system %s", systemMessage.Uuid)
		}
	} else {
		response = PeerResponse{true, "Cannot handle malformed system report", nil}
		log.WithError(err).Error("Cannot handle malformed system report")
	}

	log.Info("handleSystemReport got a response %s", response.Message)
	return response
}

func init() {
	systemIndex.HandlerFunction(SYSTEM_REPORT, handleSystemReport)
}
