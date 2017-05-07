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

package core

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"time"
)

var _ = Describe("The scheduler", func() {

	It("takes jobs and schedules them, e.g., every millisecond", func() {
		testJob := &testJob{}

		Schedule("@every 1s", testJob)

		// wait for some time and let the scheduler trigger the job several times
		time.Sleep(time.Second * 2)
		Expect(testJob.runCounter).To(BeNumerically(">=", 2))
	})

})

type testJob struct {
	runCounter uint8
}

func (t *testJob) Run() {
	t.runCounter += 1
}
