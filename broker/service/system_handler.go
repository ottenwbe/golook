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
	. "github.com/ottenwbe/golook/broker/repository"
	. "github.com/ottenwbe/golook/broker/utils"
	log "github.com/sirupsen/logrus"
)

const (
	SYSTEM_REPORT = "system_report"
)

func handleSystemReport(params interface{}) interface{} {

	var systemMessage SystemTransfer
	if err := UnmarshalB(params, &systemMessage); err == nil {
		if systemMessage.IsDeletion {
			GoLookRepository.StoreSystem(systemMessage.Uuid, systemMessage.System)
		} else {
			GoLookRepository.DelSystem(systemMessage.Uuid)
		}
	} else {
		log.WithError(err).Error("Could not handle system report")
	}
	return nil
}

func init() {
	systemIndex.HandlerFunction(SYSTEM_REPORT, handleSystemReport)
}
