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
	"github.com/ottenwbe/golook/utils"
	log "github.com/sirupsen/logrus"
)

const (
	systemReport = "system_report"
)

type (
	/*SystemService stores system information in the local registry and disseminates system information to peers.*/
	SystemService struct {
		router *router
	}
	/*SystemCallback is the base type for functions that are called by the system service.*/
	SystemCallback func(uuid string, system map[string]*golook.System)

	/*SingleSystemCallback is the base type for functions that are called by the system service.*/
	SingleSystemCallback func(uuid string, system *golook.System)

	/*SystemCallbacks is the base type for an index/map of functions that implement a 'SystemCallback'*/
	SystemCallbacks map[string]SystemCallback

	/*SingleSystemCallbacks is the base type for an index/map of functions that implement a 'SingleSystemCallback'*/
	SingleSystemCallbacks map[string]SingleSystemCallback
)

var (
	delSystemCallbacks     = SingleSystemCallbacks{}
	newSystemCallbacks     = SystemCallbacks{}
	changedSystemCallbacks = SystemCallbacks{}
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

func (s SystemService) close() {
	s.router.close()
}

/*
Run can be triggered to store and report the system's information.
*/
func (s SystemService) Run() {
	s.storeSystem(golook.GolookSystem)
	s.reportSystem()

}

func (SystemService) storeSystem(system *golook.System) {
	if system == nil {
		systemServiceLogger().Error("Ignoring nil system. This might result in your system not being reported.")
	} else {
		repo.GoLookRepository.StoreSystem(system.UUID, system)
	}
}

func (s SystemService) reportSystem() {
	var report = peerSystemReport{
		UUID:       golook.GolookSystem.UUID,
		System:     GetSystems(),
		IsDeletion: false,
	}
	response := s.broadcastSystem(report)
	systemServiceLogger().Debug(utils.MarshalSD(response))
	if response != nil {
		var systems peerSystemReport
		err := response.Unmarshal(&systems)
		if err != nil {
			systemServiceLogger().WithError(err).Error("Cannot unmarshal system!")
		} else {
			handleAddSystem(systems)
		}
	} else {
		systemServiceLogger().Error("Nil response!")
	}
}

func (s SystemService) broadcastSystem(report peerSystemReport) models.EncapsulatedValues {
	if s.router != nil {
		return s.router.BroadCast(systemReport, report)
	}
	systemServiceLogger().Error("Router is not set.")
	return nil
}

func (s SystemService) handleSystemReport(params models.EncapsulatedValues) interface{} {
	var (
		systemReport peerSystemReport
	)

	if params == nil {
		systemServiceLogger().Error("Cannot handle 'nil' system report.")
		return peerSystemReport{}
	}

	if err := params.Unmarshal(&systemReport); err != nil {
		systemServiceLogger().WithError(err).Error("Cannot handle malformed system report.")
		return peerSystemReport{}
	}

	return processSystemReport(systemReport)
}

func processSystemReport(systemReport peerSystemReport) peerSystemReport {
	if systemReport.IsDeletion {
		handleDelSystem(systemReport)
	} else {
		handleAddSystem(systemReport)
	}

	result := peerSystemReport{golook.GolookSystem.UUID, repo.GoLookRepository.GetSystems(), false}
	return result
}

func handleDelSystem(systemReport peerSystemReport) {
	s := repo.GoLookRepository.DelSystem(systemReport.UUID)
	delSystemCallbacks.call(s.UUID, s)
}

func handleAddSystem(systemReport peerSystemReport) {

	var firstTimeReport = false

	for _, s := range systemReport.System {
		_, foundSystem := repo.GoLookRepository.GetSystem(systemReport.UUID)
		firstTimeReport = firstTimeReport || foundSystem
		repo.GoLookRepository.StoreSystem(s.UUID, s)
	}

	systemServiceLogger().
		Debugf("Handle #%d systems from %s and %d callbacks.",
			len(systemReport.System),
			systemReport.UUID,
			len(changedSystemCallbacks))

	changedSystemCallbacks.call(systemReport.UUID, systemReport.System)
	if !firstTimeReport {
		newSystemCallbacks.call(systemReport.UUID, systemReport.System)
	}
}

func (s SystemService) merge(raw1 models.EncapsulatedValues, raw2 models.EncapsulatedValues) interface{} {

	var systems1 peerSystemReport
	err1 := raw1.Unmarshal(&systems1)

	var systems2 peerSystemReport
	err2 := raw2.Unmarshal(&systems2)

	if err1 != nil {
		systems1.UUID = systems2.UUID
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

/*
GetSystems return the current view on systems known the broker
*/
func GetSystems() map[string]*golook.System {
	return repo.GoLookRepository.GetSystems()
}

/*
Add a callback with a given id
*/
func (c *SystemCallbacks) Add(id string, callback SystemCallback) {
	(*c)[id] = callback
}

/*
Delete a callback with a given id
*/
func (c *SingleSystemCallbacks) Delete(id string) {
	delete(*c, id)
}

func (c *SingleSystemCallbacks) call(uuid string, system *golook.System) {
	for _, callback := range *c {
		callback(uuid, system)
	}
}

/*
Add a callback with a given id
*/
func (c *SingleSystemCallbacks) Add(id string, callback SingleSystemCallback) {
	(*c)[id] = callback
}

/*
Delete a callback with a given id
*/
func (c *SystemCallbacks) Delete(id string) {
	delete(*c, id)
}

func (c *SystemCallbacks) call(uuid string, system map[string]*golook.System) {
	for _, callback := range *c {
		callback(uuid, system)
	}
}

func systemServiceLogger() *log.Entry {
	return log.WithField("service", "system")
}
