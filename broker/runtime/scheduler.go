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

package runtime

/*
Scheduler implements a facade for a scheduler framework.
*/

import (
	"github.com/bamzi/jobrunner"
	"github.com/ottenwbe/golook/broker/utils"
	"gopkg.in/robfig/cron.v2"
)

/*
Job is a facade for the actual cron job type.
*/
type Job cron.Job

func init() {
	// Note, that jobrunner.Start() would print out to StdOut. To this end, this message is intercepted and discarded.
	utils.InterceptStdOut(func() { jobrunner.Start() })
}

/*
Schedule triggers a job regularly. When or how often it is triggered can be defined by the specification.
Examples for the specification:
"@every 5m0s"
*/
func Schedule(specification string, job Job) {
	jobrunner.Schedule(specification, cron.Job(job))
	jobrunner.Now(cron.Job(job))
}
