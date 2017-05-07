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
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	repo "github.com/ottenwbe/golook/broker/repository"
	"github.com/ottenwbe/golook/broker/routing"
	golook "github.com/ottenwbe/golook/broker/runtime/core"
)

var _ = Describe("The system service", func() {

	var (
		s          *SystemService
		mockRouter *routing.MockRouter
		sysUUID    = golook.GolookSystem.UUID
	)

	BeforeEach(func() {
		mockRouter = routing.NewMockedRouter().(*routing.MockRouter)
		s = &SystemService{router: &router{mockRouter}}
		repo.GoLookRepository.DelSystem(sysUUID)
	})

	AfterEach(func() {
		repo.GoLookRepository.DelSystem(sysUUID)
	})

	It("stores and forwards the system's information", func() {

		s.Run()

		_, isSystemInRepo := repo.GoLookRepository.GetSystem(sysUUID)
		Expect(isSystemInRepo).To(BeTrue())
		Expect(mockRouter.Visited).To(Equal(1))

	})

	It("handles valid system reports, stores the system.", func() {

		report := PeerSystemReport{Uuid: sysUUID, System: map[string]*golook.System{golook.GolookSystem.UUID: golook.GolookSystem}, IsDeletion: false}

		result := s.handleSystemReport(testEncapsulatedSystemReport{report: &report})

		_, systemInRepo := repo.GoLookRepository.GetSystem(sysUUID)
		Expect(systemInRepo).To(BeTrue())
		Expect(result).ToNot(BeNil())

	})

	It("handles valid system reports that delete systems. As a result, it removes the system, and returns a valid peer response.", func() {

		addReport := PeerSystemReport{Uuid: sysUUID, System: map[string]*golook.System{golook.GolookSystem.UUID: golook.GolookSystem}, IsDeletion: false}
		s.handleSystemReport(testEncapsulatedSystemReport{report: &addReport})
		delReport := PeerSystemReport{Uuid: sysUUID, System: map[string]*golook.System{golook.GolookSystem.UUID: golook.GolookSystem}, IsDeletion: true}
		result := s.handleSystemReport(testEncapsulatedSystemReport{report: &delReport})

		_, systemInRepo := repo.GoLookRepository.GetSystem(sysUUID)
		Expect(systemInRepo).To(BeFalse())
		Expect(result).ToNot(BeNil())

	})

	It("handles invalid system reports and returns an empty result.", func() {

		result := s.handleSystemReport(nil)

		_, systemInRepo := repo.GoLookRepository.GetSystem(sysUUID)
		Expect(systemInRepo).To(BeFalse())
		Expect(result).ToNot(BeNil())
	})

})

type testEncapsulatedSystemReport struct {
	report *PeerSystemReport
}

func (e testEncapsulatedSystemReport) Unmarshal(v interface{}) error {
	*v.(*PeerSystemReport) = *e.report
	return nil
}
