package handler

import (
	"encoding/hex"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"iris/constant"
	"iris/global"
	"iris/utils"
	"iris/web/model"
	"iris/web/service"
	"net/http"
	"strings"
	"time"
)

type IrisHandler struct {
	IrisService *service.IrisService
}

func NewIrisHandler() *IrisHandler {
	return &IrisHandler{
		IrisService: service.NewIrisService(),
	}
}

type Request struct {
	Id     string `json:"id"`
	IdType string `json:"id_type"`
	Region string `json:"region"`
}

func (h *IrisHandler) RegisterIris(ctx *gin.Context) {

	request := &Request{}
	BindJsonAndValid(ctx, request)
	if "" == strings.TrimSpace(request.Id) {
		Response(ctx, model.ResultError(http.StatusBadRequest, "ID 不能为空"))
	}
	if request.IdType == "" {
		request.IdType = constant.DefaultIdTypeHM
	}

	if request.Region == "" {
		request.Region = constant.DefaultRegionMacao
	}

	h.IrisService.RegisterIris(request.Id, request.IdType, request.Region)
	var res = "采集失败"
	for {
		select {
		case res = <-global.GVars.UartClient.OutBuffer:
			if utils.SubStr(res, 16, 2) == constant.HexCmdIrisCollect {
				errorCode := utils.SubStr(res, 22, 2)
				if errorCode == "00" {
					res = "采集成功"
					Response(ctx, model.ResultOk(res))
					return
				} else {
					res = "采集失败"
					Response(ctx, model.ResultError(http.StatusOK, res))
					return
				}

			}
		case <-time.After(15 * time.Second):
			Response(ctx, model.ResultError(http.StatusOK, res))
			return
		default:

		}
	}

	Response(ctx, model.ResultError(http.StatusOK, res))
}

func (h *IrisHandler) MatchIris(ctx *gin.Context) {
	global.GVars.UartClient.SendMsg(utils.TrimHexStr(constant.IrisDataMatch))
	var res = ""
	for {
		select {
		case res = <-global.GVars.UartClient.OutBuffer:
			if utils.SubStr(res, 16, 2) == constant.HexCmdIrisMatch {
				length := utils.SubStr(res, 18, 4)
				if length == constant.HexLengthUserID {
					bytes, _ := hex.DecodeString(strings.Trim(utils.SubStr(res, 22, 64), "0"))
					res = string(bytes)
				} else {
					res = ""
				}
				Response(ctx, model.ResultOk(res))
				return
			}
		case <-time.After(15 * time.Second):
			Response(ctx, model.ResultOk(res))
			return
		default:

		}
	}

	Response(ctx, model.ResultOk(res))
}

// DeleteAllUser /*
func (h *IrisHandler) DeleteAllUser(ctx *gin.Context) {
	data := "55 00 2A 00 00 00 00 00 44 00 20 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 E4 03"
	global.GVars.UartClient.SendMsg(utils.TrimHexStr(data))
	Response(ctx, model.ResultOk(nil))
}

// ChangeMode /**
func (h *IrisHandler) ChangeMode(ctx *gin.Context) {
	data := "55 00 0E 00 00 00 00 00 7A 00 04 00 00 00 00 DA 03"
	global.GVars.UartClient.SendMsg(utils.TrimHexStr(data))
	Response(ctx, model.ResultOk(nil))
}

// StartCycleMatch /**
// 循环识别和识别接口发送一样的命令到虹膜模组, 但需要循环发送

type Result struct {
	data string
}

func (h *IrisHandler) StartCycleMatch(ctx *gin.Context) {

	if len(global.GVars.StopCycleChan) > 0 {
		<-global.GVars.StopCycleChan
	}

	res := &Result{
		data: "",
	}
	for {
		select {
		case <-global.GVars.StopCycleChan:
			return
		default:
			log.Infof("res: %v", res)

			if res.data != "" {
				log.Infof("匹配成功！" + res.data)
				Response(ctx, model.ResultOk(res.data))
				return
			}
			global.GVars.UartClient.SendMsg(utils.TrimHexStr(constant.IrisDataMatch))
			go getMatchRes(res)
			time.Sleep(6 * time.Second)
		}
	}
	Response(ctx, model.ResultOk(res))
}

// 获取循环识别的返回结果
func getMatchRes(s *Result) {
	res := ""
	for {
		select {
		case res = <-global.GVars.UartClient.OutBuffer:

			if utils.SubStr(res, 16, 2) == constant.HexCmdIrisMatch {
				length := utils.SubStr(res, 18, 4)
				if length == constant.HexLengthUserID {
					bytes, _ := hex.DecodeString(strings.Trim(utils.SubStr(res, 22, 64), "0"))
					res = string(bytes)
					s.data = res
					log.Info("循环识别返回：" + s.data)
					return
				}
			}

		default:

		}
	}

}

func (h *IrisHandler) StopCycleMatch(ctx *gin.Context) {
	if len(global.GVars.StopCycleChan) == 0 {
		global.GVars.StopCycleChan <- struct{}{}
	}
	Response(ctx, model.ResultOk(nil))
}
