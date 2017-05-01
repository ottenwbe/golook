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

import (
	"encoding/json"
	"io"
	"net"
	"os"
	"runtime"

	"github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"
)

type System struct {
	Name string `json:"name"`
	OS   string `json:"os"`
	IP   string `json:"ip"`
	UUID string `json:"uuid"`
}

var (
	GolookSystem *System = NewSystem()
)

func NewSystem() *System {
	return &System{
		Name: getName(),
		OS:   getOS(),
		IP:   getIP(),
		UUID: getUuid(),
	}
}

func EncodeSystem(sys *System) string {
	b, err := json.Marshal(sys)
	return toString(b, err)
}

func EncodeSystems(sys []*System) string {
	b, err := json.Marshal(sys)
	return toString(b, err)
}

func toString(b []byte, err error) string {
	if err != nil {
		log.WithError(err).Error("Error when marshalling system")
		return "{}"
	} else {
		return string(b)
	}
}

func DecodeSystems(reader io.Reader) ([]*System, error) {
	systems := make([]*System, 0)
	err := json.NewDecoder(reader).Decode(&systems)
	return systems, err
}

func DecodeSystem(sysReader io.Reader) (System, error) {
	var sys System
	err := json.NewDecoder(sysReader).Decode(&sys)
	return sys, err
}

func getUuid() string {
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

//https://play.golang.org/p/BDt3qEQ_2H
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
