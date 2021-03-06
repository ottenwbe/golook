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
	"github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"
	"net"
	"os"
	"runtime"
)

/*
System represents the runtime system (OS, ...) of the broker.
*/
type System struct {
	Name string `json:"name"`
	OS   string `json:"os"`
	IP   string `json:"ip"`
	UUID string `json:"uuid"`
}

var (
	/*GolookSystem is a reference to the information about the system on which this app is running*/
	GolookSystem = NewSystem()
)

/*
NewSystem is the factory function for System and returns by default all information about the system on which this golook app is running
*/
func NewSystem() *System {
	return &System{
		Name: getName(),
		OS:   getOS(),
		IP:   getIP(),
		UUID: getUUID(),
	}
}

func getUUID() string {
	return uuid.NewV5(uuid.NamespaceURL, getIP()).String()
}

func getName() string {
	hostName, err := os.Hostname()
	logError(err)
	return hostName
}

func getOS() string {
	return runtime.GOOS
}

//see: https://play.golang.org/p/BDt3qEQ_2H
func getIP() string {
	ifaces, err := net.Interfaces()
	if err != nil {
		return ""
	}
	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 {
			continue // interface down
		}
		if iface.Flags&net.FlagLoopback != 0 {
			continue // loopback interface
		}
		addrs, err := iface.Addrs()
		if err != nil {
			return ""
		}
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if ip == nil || ip.IsLoopback() {
				continue
			}
			ip = ip.To4()
			if ip == nil {
				continue // not an ipv4 address
			}
			return ip.String()
		}
	}
	return "" //errors.New("No connection detected")
}

func logError(err error) {
	if err != nil {
		log.WithError(err).Error("Error when instantiating System")
	}
}
