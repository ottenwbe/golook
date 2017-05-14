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
	"github.com/ottenwbe/golook/broker/models"
	log "github.com/sirupsen/logrus"
	"os"
	"path/filepath"
	"time"
)

type testReporter struct {
	monitorReports int
}

func (t *testReporter) reportHandler(file string) {
	t.monitorReports++
}

var _ = Describe("The file monitor", func() {

	const (
		testFile  = "test.txt"
		testFile2 = "test2.txt"
	)

	var (
		fm *FileMonitor
		tr *testReporter
	)

	BeforeEach(func() {
		tr = &testReporter{}
		fm = &FileMonitor{
			Reporter: tr.reportHandler,
		}
		fm.Open(map[string]map[string]*models.File{})
	})

	AfterEach(func() {
		fm.Close()
		fm = nil
	})

	It("is triggered by adding and removing a file in a monitored folder.", func() {

		// Monitor the current directory for changes
		currentDirectory, err := filepath.Abs(".")
		if err != nil {
			log.Fatal(err)
		}

		fm.Monitor(currentDirectory)

		f, err := os.Create(testFile)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()

		err = os.Remove(testFile)
		if err != nil {
			log.Fatal(err)
		}

		// wait for both events, or wait for 1 second to ensure that the test eventually stops
		numEvents := 0
		stopTime := time.Now().Add(time.Second)
		for numEvents < 1 && time.Now().Before(stopTime) {
			numEvents = tr.monitorReports
		}

		Expect(fm.watcher).ToNot(BeNil())
		// Due to timings, the test might be flaky and miss the create event
		// Therefore, only check for at least one event that
		Expect(numEvents >= 1).To(BeTrue())

	})

	It("is triggered by adding and removing a file, if it is monitored.", func() {
		f, err := os.Create(testFile2)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()
		defer os.Remove(testFile2)

		fm.Monitor(testFile2)

		time.Sleep(time.Millisecond * 100)

		f.WriteString("test write")

		// wait for the write event, or wait for 1 second to ensure that the test stops eventually
		numEvents := 0
		stopTime := time.Now().Add(time.Second)
		for numEvents < 1 && time.Now().Before(stopTime) {
			numEvents = tr.monitorReports
		}

		Expect(numEvents).To(BeNumerically(">", 0))

	})
})

var _ = Describe("The file monitor's initialization", func() {

	const testFile = "test-init.txt"

	var (
		fm *FileMonitor
		tr *testReporter
	)

	BeforeEach(func() {
		tr = &testReporter{}
		fm = &FileMonitor{
			Reporter: tr.reportHandler,
		}
	})

	AfterEach(func() {
	})

	It("ensures that a file watcher is created after the initialization", func() {
		fm.Open(map[string]map[string]*models.File{})
		defer fm.Close()

		Expect(fm.watcher).ToNot(BeNil())
	})

	It("ensures that a file watcher is created after the initializon, even if the defaults for monitored files are nil.", func() {
		fm.Open(nil)
		defer fm.Close()

		Expect(fm.watcher).ToNot(BeNil())
	})

	It("monitors files that are given as default.", func() {

		f, err := os.Create(testFile)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()
		defer os.Remove(testFile)

		monitoredFile, _ := models.NewFile(testFile)

		fm.Open(map[string]map[string]*models.File{".": {monitoredFile.Name: monitoredFile}})
		defer fm.Close()

		f.WriteString("test write")

		// wait for the write event, or wait for 1 second to ensure that the test stops eventually
		numEvents := 0
		stopTime := time.Now().Add(time.Second)
		for numEvents < 1 && time.Now().Before(stopTime) {
			numEvents = tr.monitorReports
		}

		Expect(numEvents).To(BeNumerically(">", 0))
	})

})
