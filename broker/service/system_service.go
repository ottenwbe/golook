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
	repo "github.com/ottenwbe/golook/broker/repository"
	"github.com/ottenwbe/golook/broker/routing"
	golook "github.com/ottenwbe/golook/broker/runtime/core"
	log "github.com/sirupsen/logrus"
)

const (
	systemReport = "system_report"
)

type (
	/*
		SystemService stores system information in the local registry and disseminates system information to peers.
	*/
	SystemService struct {
		router *router
	}
	NewSystemCallback  func(uuid string, system map[string]*golook.System)
	NewSystemCallbacks map[string]NewSystemCallback
)

var (
	newSystemCallbacks = NewSystemCallbacks{}
)

func newSystemService() *SystemService {
	s := &SystemService{}
	s.router = newRouter("broadcast_system", routing.BroadcastRouter)
	s.router.AddHandler(
		systemReport,
		routing.NewHandler(s.handleSystemReport, s.merge),
	)
	return s
}

func (s SystemService) merge(raw1 models.EncapsulatedValues, raw2 models.EncapsulatedValues) interface{} {

	var systems1 PeerSystemReport
	err1 := raw1.Unmarshal(&systems1)

	var systems2 PeerSystemReport
	err2 := raw2.Unmarshal(&systems2)

	if err1 != nil {
		systems1.Uuid = systems2.Uuid
		systems1.System = make(map[string]*golook.System)
	}

	if err2 != nil {
		systems2.System = make(map[string]*golook.System)
	}

	for k, v := range systems2.System {
		systems1.System[k] = v
	}

	return systems1
}

func (s SystemService) close() {
	s.router.close()
}

/*
Run can be triggered to store and report the system's information.
*/
func (s SystemService) Run() {
	s.storeSystem(golook.GolookSystem)

	var report = PeerSystemReport{
		Uuid:       golook.GolookSystem.UUID,
		System:     repo.GoLookRepository.GetSystems(),
		IsDeletion: false,
	}
	response := s.broadcastSystem(report)
	if response != nil {
		var systems PeerSystemReport
		response.Unmarshal(systems)
		handleNewSystem(systems)
	} else {
		log.Error("Nil response!")
	}
}

func (SystemService) storeSystem(system *golook.System) {
	if system != nil {
		repo.GoLookRepository.StoreSystem(system.UUID, system)
	} else {
		log.Error("Ignoring nil system when storing it.")
	}
}

func (s SystemService) broadcastSystem(report PeerSystemReport) models.EncapsulatedValues {
	if s.router != nil {
		return s.router.BroadCast(systemReport, report)
	}
	log.WithField("service", "system").Error("router is not set!")
	return nil
}

func (s SystemService) handleSystemReport(params models.EncapsulatedValues) interface{} {
	var (
		systemReport PeerSystemReport
	)

	if params == nil {
		log.Error("Cannot handle 'nil' system report")
		return PeerSystemReport{}
	}

	if err := params.Unmarshal(&systemReport); err != nil {
		log.WithError(err).Error("Cannot handle malformed system report")
		return PeerSystemReport{}
	}

	return processSystemReport(systemReport)
}

func processSystemReport(systemReport PeerSystemReport) PeerSystemReport {
	var err error
	if systemReport.IsDeletion {
		repo.GoLookRepository.DelSystem(systemReport.Uuid)
	} else {
		err = handleNewSystem(systemReport)
	}

	if err != nil {
		log.WithError(err).Info("Error processing system report.")
	}

	result := PeerSystemReport{golook.GolookSystem.UUID, repo.GoLookRepository.GetSystems(), false}
	return result
}

func handleNewSystem(systemMessage PeerSystemReport) error {

	_, found := repo.GoLookRepository.GetSystem(systemMessage.Uuid)

	for _, s := range systemMessage.System {
		repo.GoLookRepository.StoreSystem(s.UUID, s)
	}

	log.WithField("found", found).Infof("Got #%d systems at %s from %s and %d callbacks",
		len(systemMessage.System),
		golook.GolookSystem.UUID,
		systemMessage.Uuid,
		len(systemMessage.System))

	if !found {
		newSystemCallbacks.call(systemMessage.Uuid, systemMessage.System)
	}
	return nil
}

func GetSystems() map[string]*golook.System {
	return repo.GoLookRepository.GetSystems()
}

func (c *NewSystemCallbacks) Add(id string, callback NewSystemCallback) {
	(*c)[id] = callback
}

func (c *NewSystemCallbacks) Delete(id string) {
	delete(*c, id)
}

func (c *NewSystemCallbacks) call(uuid string, system map[string]*golook.System) {
	for _, callback := range *c {
		callback(uuid, system)
	}
}
