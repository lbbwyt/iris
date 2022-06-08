package handler

import (
	"encoding/hex"
	getty "github.com/AlexStocks/getty/transport"
	log "github.com/sirupsen/logrus"
	"iris/constant"
	"iris/global"
	"iris/utils"
	"strings"
)

type IrisMessageHandler struct{}

func NewEchoMessageHandler() *IrisMessageHandler {
	return &IrisMessageHandler{}
}

func (h *IrisMessageHandler) OnOpen(session getty.Session) error {
	log.Info("client connected:" + session.RemoteAddr())
	global.GVars.Conn = session
	return nil
}

func (h *IrisMessageHandler) OnError(session getty.Session, err error) {
	log.Errorf("session{%s} got error{%v}.", session.Stat(), err)
}

func (h *IrisMessageHandler) OnClose(session getty.Session) {
	log.Infof("session{%s} is closing......", session.Stat())
}

func (h *IrisMessageHandler) OnMessage(session getty.Session, pkg interface{}) {
	log.Infof("get echo package{%s}", pkg)
	global.GVars.UartClient.SendMsg(pkg.(string))

	for {
		select {
		case res := <-global.GVars.UartClient.OutBuffer:
			if utils.SubStr(res, 16, 2) == constant.HexCmdIrisMatch {
				length := utils.SubStr(res, 18, 4)
				if length == constant.HexLengthUserID {
					bytes, _ := hex.DecodeString(strings.Trim(utils.SubStr(res, 22, 64), "0"))
					res = string(bytes)
					log.Infof("返回tcp客户端：" + res)
					global.GVars.Conn.WriteBytes(bytes)
				}
				return

			}
		}
	}

}

func (h *IrisMessageHandler) OnCron(session getty.Session) {
	//心跳处理
	log.Info("心跳处理")
}
