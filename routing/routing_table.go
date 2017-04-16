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
package routing

import (
	. "github.com/ottenwbe/golook/rpc_client"
)

var (
	GolookClient LookupClient
	//Downlink     []LookupClient //TODO store downlink information
)

func ConfigLookClient(host string, port int) {
	GolookClient = NewLookClient(host, port)
}

func AccessMockedGolookClient() *MockGolookClient {
	return GolookClient.(*MockGolookClient)
}

func RunWithMockedGolookClient(mockedFunction func()) {
	RunWithMockedGolookClientF(mockedFunction, "", "")
}

func RunWithMockedGolookClientF(mockedFunction func(), fileName string, folderName string) {
	//ensure that the GolookClient is reset after the function's execution
	defer func(reset LookupClient) {
		GolookClient = reset
	}(GolookClient)

	//create a mock rpc_client
	GolookClient = &MockGolookClient{
		VisitDoPostFile:  false,
		VisitDoPutFiles:  false,
		VisitDoGetFiles:  false,
		VisitDoPostFiles: false,
		FileName:         fileName,
		FolderName:       folderName,
	}

	mockedFunction()
}

/*var _ = Describe(" LookupClient ", func() {
	It("should be configured during construction with host and port", func() {
		ConfigLookClient("do.test", 8123)
		Expect(GolookClient.(*LookupClientData).serverUrl).To(ContainSubstring("do.test"))
		Expect(GolookClient.(*LookupClientData).serverUrl).To(ContainSubstring("8123"))
	})
})

var _ = Describe(" MockGoLookClient ", func() {
	It("should be set when a function is run in context of 'RunWithMockedGolookClient'", func() {
		RunWithMockedGolookClient(func() {
			Expect(testConvertMockGolookClient).ToNot(Panic())
		})
	})

	It("should be set when a function is run in context of 'RunWithMockedGolookClientF'", func() {
		RunWithMockedGolookClientF(func() {
			Expect(testConvertMockGolookClient).ToNot(Panic())
			Expect(GolookClient.(*MockGolookClient).fileName).To(Equal("fileName.txt"))
			Expect(GolookClient.(*MockGolookClient).folderName).To(Equal("folder"))
		}, "fileName.txt", "folder")
	})
})

func testConvertMockGolookClient() {
	_ = GolookClient.(*MockGolookClient)
}

*/
