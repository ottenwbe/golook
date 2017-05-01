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
	golook "github.com/ottenwbe/golook/broker/runtime"
)

var _ = Describe("The system handler", func() {

	It("handles valid system reports and returns a peer response", func() {

		sysUUID := golook.GolookSystem.UUID
		repo.GoLookRepository.DelSystem(sysUUID)
		report := PeerSystemReport{Uuid: sysUUID, System: golook.GolookSystem, IsDeletion: false}

		result := handleSystemReport(testEncapsulatedSystemReport{report: &report})

		_, systemInRepo := repo.GoLookRepository.GetSystem(sysUUID)
		Expect(systemInRepo).To(BeTrue())
		Expect(result).ToNot(BeNil())
		Expect(result.(PeerResponse).Error).To(BeFalse())

	})
})

type testEncapsulatedSystemReport struct {
	report *PeerSystemReport
}

func (e testEncapsulatedSystemReport) Unmarshal(v interface{}) error {
	*v.(*PeerSystemReport) = *e.report
	return nil
}
