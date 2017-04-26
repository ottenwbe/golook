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

//
import (
	"encoding/json"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/ottenwbe/golook/broker/models"
	. "github.com/ottenwbe/golook/broker/repository"
	"github.com/ottenwbe/golook/broker/runtime"
)

var _ = Describe("The file handler", func() {

	BeforeEach(func() {
		//reset repo ... for each test
		GoLookRepository = NewRepository()
	})

	AfterEach(func() {
		//reset repo ... would actually be sufficient after the last test
		GoLookRepository = NewRepository()
	})

	It("stores file reports in the golook repository (when the system is known to the repo)", func() {
		storedSys := &models.SystemFiles{System: runtime.GolookSystem, Files: nil}
		GoLookRepository.StoreSystem(runtime.GolookSystem.UUID, storedSys)
		f, _ := models.NewFile("file_handler_test.go")
		fReport, _ := json.Marshal(PeerFileReport{Files: map[string]*models.File{f.Name: f}, Replace: false, System: runtime.GolookSystem.UUID})
		handleFileReport(fReport)
		mapRepo := AccessMapRepository()
		Expect(len((*mapRepo)[runtime.GolookSystem.UUID].Files)).To(Equal(1))
	})

	It("rejects invalid reports)", func() {
		//GoLookRepository.StoreSystem(runtime.GolookSystem.UUID, runtime.GolookSystem)
		fReport := []byte("")
		handleFileReport(fReport)
		Expect(len(*GoLookRepository.(*MapRepository))).To(BeZero())
	})
})