package app

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
)

const (
	GOLOOK_NAME = "golook"
	VERSION     = "v0.1.0-dev"
)

type Info struct {
	App     string  `json:"app"`
	Version string  `json:"version"`
	System  *System `json:"system"`
}

func NewAppInfo() *Info {
	return &Info{
		App:     GOLOOK_NAME,
		Version: VERSION,
		System:  NewSystem(),
	}
}

func EncodeAppInfo(info *Info) string {
	b, err := json.Marshal(info)
	if err != nil {
		logrus.WithError(err).Error("Could not encode app info.")
		return "{}"
	}
	return string(b)
}
