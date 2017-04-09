package control

import (
	. "github.com/ottenwbe/golook/routing"
	. "github.com/ottenwbe/golook/utils"
)

type SystemController struct {
	uplink string
	system *System
}

func NewSystemController(uplink string) *SystemController {
	return &SystemController{
		uplink: uplink,
		system: NewSystem(),
	}
}

func (sc *SystemController) Connect() {
	GolookClient.DoPutSystem(sc.system)
}

func (sc *SystemController) ConnectWith(uplinkHost string) {
	GolookClient.DoPutSystem(sc.system)
}

func (sc *SystemController) Disconnect(uplinkHost string) {

}

func (sc *SystemController) DisconnectAll() {

}
