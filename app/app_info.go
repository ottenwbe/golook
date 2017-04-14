package app

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
)

const (
	APP_NAME = "golook"
	VERSION  = "v0.1.0-dev"
)

type AppInfo struct {
	App     string  `json:"app"`
	Version string  `json:"version"`
	System  *System `json:"system"`
}

func NewAppInfo() *AppInfo {
	return &AppInfo{
		App:     APP_NAME,
		Version: VERSION,
		System:  NewSystem(),
	}
}

func EncodeAppInfo(info *AppInfo) string {
	b, err := json.Marshal(info)
	if err != nil {
		logrus.WithError(err).Error("Could not encode app info.")
		return "{}"
	}
	return string(b)
}
