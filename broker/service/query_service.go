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
	"github.com/ottenwbe/golook/broker/repository"
	"github.com/ottenwbe/golook/broker/utils"
	log "github.com/sirupsen/logrus"
)

const (
	MockQueries  = "mock"
	LocalQueries = "local"
	BCastQueries = "broadcast"
)

type (
	queryService interface {
		query(searchString string) (interface{}, error)
	}
	localQueryService     struct{}
	broadcastQueryService struct {
		router *Router
	}
	MockQueryService struct {
		SearchString string
	}
	FileQueryData map[string][]*models.File
)

func newQueryService(queryType string, router *Router) queryService {
	switch queryType {
	case MockQueries:
		return &MockQueryService{}
	case BCastQueries:
		return &broadcastQueryService{router: router}
	default:
		return &localQueryService{}

	}
}

func (*localQueryService) query(searchString string) (interface{}, error) {
	fq := PeerFileQuery{SearchString: searchString}
	return processFileQuery(fq), nil
}

func (qs *broadcastQueryService) query(searchString string) (interface{}, error) {
	fq := PeerFileQuery{SearchString: searchString}
	queryResult := qs.router.BroadCast(fileQuery, fq)

	var response FileQueryData
	err := utils.Unmarshal(queryResult, &response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (mock *MockQueryService) query(searchString string) (interface{}, error) {
	mock.SearchString = searchString
	return "{}", nil
}

const (
	fileQuery = "file query"
)

func mergeFileQuery(v1 models.EncapsulatedValues, v2 models.EncapsulatedValues) interface{} {

	rawData1, rawData2 := getRawData(v1, v2)
	mergedData := mergeData(rawData1, rawData2)

	return mergedData
}

func mergeData(rawData1 FileQueryData, rawData2 FileQueryData) FileQueryData {
	result := rawData1
	for k, v := range rawData2 {
		result[k] = v
	}
	return result
}

func getRawData(v1 models.EncapsulatedValues, v2 models.EncapsulatedValues) (FileQueryData, FileQueryData) {
	var (
		response1, response2 map[string][]*models.File
	)

	errV1 := v1.Unmarshal(&response1)
	errV2 := v2.Unmarshal(&response2)

	if errV1 != nil {
		response1 = map[string][]*models.File{}
	}
	if errV2 == nil {
		response2 = map[string][]*models.File{}
	}

	return response1, response2
}

func handleFileQuery(params models.EncapsulatedValues) interface{} {

	var (
		queryMessage PeerFileQuery
		response     FileQueryData
	)

	err := params.Unmarshal(&queryMessage)

	if err == nil {
		response = processFileQuery(queryMessage)
	} else {
		log.WithError(err).Error("Could not handle file query")
		response = FileQueryData{}
	}

	return response
}

func processFileQuery(systemMessage PeerFileQuery) FileQueryData {
	result := repositories.GoLookRepository.FindSystemAndFiles(systemMessage.SearchString)
	return FileQueryData(result)
}
