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

	"github.com/sirupsen/logrus"
	"os"
	"path/filepath"
	"time"
)

type testReporter struct {
	monitorReports int
}

func (t *testReporter) reportHandler(file string) {
	t.monitorReports += 1
}

var _ = Describe("The file monitor", func() {

	var (
		fm *FileMonitor
		tr *testReporter
	)

	BeforeEach(func() {
		tr = &testReporter{}
		fm = &FileMonitor{
			reporter: tr.reportHandler,
		}
		fm.Start()
	})

	AfterEach(func() {
		fm.Close()
		fm = nil
	})

	Context("initialization", func() {
		It("is running after the initialization", func() {
			Expect(fm.watcher).ToNot(BeNil())
		})
	})

	Context("monitoring", func() {

		const TEST_FILE = "test.txt"

		It("is triggered by adding and removing a file", func() {

			// Monitor the current directory for changes
			currentDirectory, _ := filepath.Abs(".")
			fm.Monitor(currentDirectory)

			_, err := os.Create(TEST_FILE) // For read access.
			if err != nil {
				logrus.Fatal(err)
			}

			err = os.Remove(TEST_FILE) // For read access.
			if err != nil {
				logrus.Fatal(err)
			}

			// wait for both events, or wait for 1 second to ensure that the test eventually stops
			numEvents := 0
			elapsed := time.Now().Second()
			for numEvents < 1 && time.Now().Second()-elapsed < 1 {
				numEvents = tr.monitorReports
			}

			Expect(fm.watcher).ToNot(BeNil())
			// Due to timings, the test might be flaky and miss the create event
			// Therefore, only check for at least one event that
			Expect(numEvents >= 1).To(BeTrue())

		})
	})
})
