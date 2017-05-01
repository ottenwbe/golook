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

	log "github.com/sirupsen/logrus"

	"github.com/ottenwbe/golook/broker/models"
	repo "github.com/ottenwbe/golook/broker/repository"
	"github.com/ottenwbe/golook/broker/runtime"
)

const (
	SYSTEM_REPORT = "system_report"
)

func handleSystemReport(params models.EncapsulatedValues) interface{} {
	var (
		systemReport PeerSystemReport
	)

	if err := params.Unmarshal(&systemReport); err != nil {
		log.WithError(err).Error("Cannot handle malformed system report")
		return PeerResponse{true, "Cannot handle malformed system report", nil}
	}

	return processSystemReport(systemReport)
}

func processSystemReport(systemReport PeerSystemReport) PeerResponse {
	var err error
	if systemReport.IsDeletion {
		repo.GoLookRepository.DelSystem(systemReport.Uuid)
	} else {
		err = handleNewSystem(systemReport)
	}

	if err != nil {
		return PeerResponse{true, fmt.Sprintf("%s", err), nil}
	}
	return PeerResponse{false, fmt.Sprintf("Processed request for system %s", systemReport.Uuid), nil}

}

func handleNewSystem(systemMessage PeerSystemReport) error {
	_, found := repo.GoLookRepository.GetSystem(systemMessage.System.UUID)
	if !found {
		newSystemCallbacks.Call(systemMessage.Uuid, systemMessage.System)
	}
	repo.GoLookRepository.StoreSystem(systemMessage.Uuid, systemMessage.System)
	return nil
}

type NewSystemCallbacks map[string]func(uuid string, system *runtime.System)

var newSystemCallbacks = NewSystemCallbacks{}

func (c *NewSystemCallbacks) Add(id string, callback func(uuid string, system *runtime.System)) {
	(*c)[id] = callback
}

func (c *NewSystemCallbacks) Delete(id string) {
	delete(*c, id)
}

func (c *NewSystemCallbacks) Call(uuid string, system *runtime.System) {
	for _, callback := range *c {
		callback(uuid, system)
	}
}
