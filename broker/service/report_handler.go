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
	"github.com/ottenwbe/golook/broker/models"
	. "github.com/ottenwbe/golook/broker/repository"
	log "github.com/sirupsen/logrus"
)

const (
	fileReport = "file_report"
)

func handleFileReport(params models.EncapsulatedValues) interface{} {
	var (
		fileMessage PeerFileReport
		response    PeerResponse
	)

	if err := params.Unmarshal(&fileMessage); err == nil {
		response = processFileReport(&fileMessage)
	} else {
		log.WithError(err).Error("Cannot handle file report.")
		response = PeerResponse{Error: true, Message: "Malformed file report", Data: nil}
	}

	return response
}

func processFileReport(fileMessage *PeerFileReport) (response PeerResponse) {
	response.Error = !GoLookRepository.UpdateFiles(fileMessage.System, fileMessage.Files)
	response.Message = fmt.Sprintf("Processed file report for system %s", fileMessage.System)
	return
}
