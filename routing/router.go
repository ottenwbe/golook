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
	"github.com/geofffranks/spruce/log"
	. "github.com/ottenwbe/golook/rpc_client"
	"github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
)

type Key struct {
	id uuid.UUID
}

func NilKey() Key {
	return Key{
		id: uuid.Nil,
	}
}

func NewKey(name string) Key {
	return Key{
		id: uuid.NewV5(uuid.Nil, name),
	}
}

func NewKeyN(namespace uuid.UUID, name string) Key {
	return Key{
		id: uuid.NewV5(namespace, name),
	}
}

var (
	routeTable   DefaultRouteTable = DefaultRouteTable{}
	handlerTable HandlerTable      = HandlerTable{}
)

type HandlerTable map[string]func(params interface{})

type RouteTable interface {
	get(key Key) (LookupClient, bool)
	put(key Key, client LookupClient)
	this() LookupClient
	predecessor() LookupClient
	successor() LookupClient
}

type DefaultRouteTable struct {
	tbl               map[Key]LookupClient
	thisClient        LookupClient
	predecessorClient LookupClient
	successorClient   LookupClient
}

func (rt DefaultRouteTable) this() LookupClient {
	return rt.thisClient
}

func (rt DefaultRouteTable) predecessor() LookupClient {
	return rt.predecessorClient
}

func (rt DefaultRouteTable) successor() LookupClient {
	return rt.successorClient
}

func (rt DefaultRouteTable) get(key Key) (LookupClient, bool) {
	client, ok := rt.tbl[key]
	return client, ok
}

func (rt DefaultRouteTable) put(key Key, client LookupClient) {
	rt.tbl[key] = client
}

type Router interface {
	route(key Key, method string, params interface{})
	handle(key Key, method string, params interface{})
}

type BroadcastRouter struct {
}

func (BroadcastRouter) route(_ Key, method string, message interface{}) {
	for _, client := range routeTable.tbl {
		client.Call(method, message)
		// Make the call
		log.DEBUG("Route message to client: %s", client)
	}
}

func (BroadcastRouter) handle(_ Key, method string, message interface{}) {
	if function, ok := handlerTable[method]; ok {
		function(message)
	} else {
		logrus.Errorf("Handler for function %s not found", method)
	}
}

//
//type LookupHandler interface {
//	handleQueryAllSystemsForFile(fileName string) (files map[string]*System, err error)
//	handleQueryFiles(systemName string) (files map[string]File, err error)
//
//	handleReportFile(system string, filePath string) error
//	handleReportFileR(system string, filePath string) error
//	handleReportFolderR(system string, folderPath string) error
//	handleReportFolder(system string, folderPath string) error
//
//	handleFileDeletion(system string, filePath string) error
//	handleFolderDeletion(system string, filePath string) error
//}
//
//type ReportRouter interface {
//	handleSystemReport(system string) error
//	handleSystemDeletion(system string) error
//}
//
//type SystemRouter interface {
//	handleSystemReport(system string) error
//	handleSystemDeletion(system string) error
//}
