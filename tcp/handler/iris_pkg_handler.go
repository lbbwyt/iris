package handler

import (
	"encoding/hex"
	getty "github.com/AlexStocks/getty/transport"
	log "github.com/sirupsen/logrus"
)

type IrisPackageHandler struct{}

func NewEchoPackageHandler() *IrisPackageHandler {
	return &IrisPackageHandler{}
}

func (h *IrisPackageHandler) Read(ss getty.Session, data []byte) (interface{}, int, error) {

	log.Info("receive:" + string(data))
	return string(data), len(data), nil
}

func (h *IrisPackageHandler) Write(ss getty.Session, pkg interface{}) ([]byte, error) {
	log.Info("send:" + pkg.(string))
	return hex.DecodeString(pkg.(string))
}
