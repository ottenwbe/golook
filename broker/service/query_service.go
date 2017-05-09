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
	"github.com/ottenwbe/golook/broker/routing"
	"github.com/ottenwbe/golook/utils"
	log "github.com/sirupsen/logrus"
)

const (
	mockQueries  = "mock"
	localQueries = "local"
	bCastQueries = "broadcast"
)

type (
	queryService interface {
		query(searchString string) (interface{}, error)
	}
	localQueryService     struct{}
	broadcastQueryService struct {
		router *router
	}
	/*
		MockQueryService implements a mock interface for queries
	*/
	MockQueryService struct {
		SearchString string
	}
	fileQueryData map[string][]*models.File
)

func newQueryService(queryType string, router *router) queryService {
	switch queryType {
	case mockQueries:
		return &MockQueryService{}
	case bCastQueries:
		if router != nil {
			router.AddHandler(
				fileQuery,
				routing.NewHandler(
					handleFileQuery,
					mergeFileQuery,
				),
			)
		}
		return &broadcastQueryService{router: router}
	default:
		if router != nil {
			router.AddHandler(
				fileQuery,
				routing.NewHandler(
					handleFileQuery,
					mergeFileQuery,
				),
			)
		}
		return &localQueryService{}
	}
}

func (*localQueryService) query(searchString string) (interface{}, error) {
	fq := peerFileQuery{SearchString: searchString}
	return processFileQuery(fq), nil
}

func (qs *broadcastQueryService) query(searchString string) (interface{}, error) {
	fq := peerFileQuery{SearchString: searchString}
	queryResult := qs.router.BroadCast(fileQuery, fq)

	var response fileQueryData
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

func mergeData(rawData1 fileQueryData, rawData2 fileQueryData) fileQueryData {
	result := rawData1
	for k, v := range rawData2 {
		result[k] = v
	}
	return result
}

func getRawData(rawData1 models.EncapsulatedValues, rawData2 models.EncapsulatedValues) (fileQueryData, fileQueryData) {
	var (
		response1, response2 fileQueryData
	)

	errV1 := rawData1.Unmarshal(&response1)
	errV2 := rawData2.Unmarshal(&response2)

	if errV1 != nil {
		response1 = fileQueryData{}
	}
	if errV2 != nil {
		response2 = fileQueryData{}
	}

	return response1, response2
}

func handleFileQuery(params models.EncapsulatedValues) interface{} {

	var (
		queryMessage peerFileQuery
		response     fileQueryData
	)

	err := params.Unmarshal(&queryMessage)

	if err == nil {
		response = processFileQuery(queryMessage)
	} else {
		log.WithError(err).Error("Could not handle file query")
		response = fileQueryData{}
	}

	return response
}

func processFileQuery(systemMessage peerFileQuery) fileQueryData {
	result := repositories.GoLookRepository.FindSystemAndFiles(systemMessage.SearchString)
	return fileQueryData(result)
}
