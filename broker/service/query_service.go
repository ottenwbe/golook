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
	"github.com/ottenwbe/golook/broker/utils"
)

const (
	MockQueries  = "mock"
	LocalQueries = "local"
	BCastQueries = "broadcast"
)

type (
	queryService interface {
		MakeFileQuery(searchString string) (interface{}, error)
	}
	localQueryService     struct{}
	broadcastQueryService struct{}
	MockQueryService      struct {
		SearchString string
	}
)

func newQueryService(queryType string) queryService {
	switch queryType {
	case MockQueries:
		return &MockQueryService{}
	case BCastQueries:
		return &broadcastQueryService{}
	default:
		return &localQueryService{}

	}
}

func (*localQueryService) MakeFileQuery(searchString string) (interface{}, error) {
	fq := PeerFileQuery{SearchString: searchString}
	return string(processFileQuery(fq).Data), nil
}

func (*broadcastQueryService) MakeFileQuery(searchString string) (interface{}, error) {
	fq := PeerFileQuery{SearchString: searchString}
	queryResult := broadCastRouter.BroadCast(FILE_QUERY, fq)

	var response PeerResponse
	err := utils.Unmarshal(queryResult, &response)
	if err != nil {
		return nil, err
	}

	//TODO: response error

	return response.Data, nil
}

func (mock *MockQueryService) MakeFileQuery(searchString string) (interface{}, error) {
	mock.SearchString = searchString
	return "{}", nil
}
