package service

import (
	"encoding/hex"
	"fmt"
	log "github.com/sirupsen/logrus"
	"iris/global"
	"iris/utils"
)

type IrisService struct {
}

func NewIrisService() *IrisService {
	return &IrisService{}
}

func (s *IrisService) RegisterIris(id, idType, region string) error {

	idStr := fmt.Sprintf("%s%s%s", region, idType, id)
	log.Info(idStr)
	prefix := "55002A0000000000430020"
	suffix := "03"
	hexData := utils.HexStrFill(hex.EncodeToString(utils.StringToBytes(idStr)), 64)
	checkSum, err := utils.XrCheckSum(prefix + hexData)
	if err != nil {
		return err
	}
	data := prefix + hexData + checkSum + suffix
	log.Info("发送注册数据：" + data)
	global.GVars.UartClient.SendMsg(data)
	return nil
}
