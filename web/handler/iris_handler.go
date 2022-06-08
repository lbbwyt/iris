package handler

import (
	"encoding/hex"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"iris/constant"
	"iris/global"
	"iris/utils"
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

func (h *IrisHandler) RegisterIris(ctx *gin.Context) {

	id := ctx.Query("id")

	if "" == strings.TrimSpace(id) {
		ResponseBadRequest(ctx, "ID 不能为空")
	}
	idType := ctx.Query("id_type")
	region := ctx.Query("region")
	if idType == "" {
		idType = constant.DefaultIdTypeHM
	}
	if region == "" {
		region = constant.DefaultRegionMacao
	}

	h.IrisService.RegisterIris(id, idType, region)
	var res = "采集失败"
	for {
		select {
		case res = <-global.GVars.UartClient.OutBuffer:
			if utils.SubStr(res, 16, 2) == constant.HexCmdIrisCollect {
				errorCode := utils.SubStr(res, 22, 2)
				if errorCode == "00" {
					res = "采集成功"
					ResponseJSON(ctx, http.StatusOK, res)
				} else {
					res = "采集失败"
					ResponseServerError(ctx, res)
				}
				return
			}
		case <-time.After(15 * time.Second):
			ResponseServerError(ctx, res)
			return
		}
	}

	ResponseServerError(ctx, res)
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
				ResponseJSON(ctx, http.StatusOK, res)
				return
			}
		case <-time.After(15 * time.Second):
			ResponseJSON(ctx, http.StatusOK, res)
			return
		}
	}

	ResponseJSON(ctx, http.StatusOK, res)
}

// DeleteAllUser /*
func (h *IrisHandler) DeleteAllUser(context *gin.Context) {
	data := "55 00 2A 00 00 00 00 00 44 00 20 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 E4 03"
	global.GVars.UartClient.SendMsg(utils.TrimHexStr(data))
	ResponseOK(context)
}

// ChangeMode /**
func (h *IrisHandler) ChangeMode(context *gin.Context) {
	data := "55 00 0E 00 00 00 00 00 7A 00 04 00 00 00 00 DA 03"
	global.GVars.UartClient.SendMsg(utils.TrimHexStr(data))
	ResponseOK(context)
}

// StartCycleMatch /**

// 循环识别和识别接口发送一样的命令到虹膜模组, 但需要循环发送

type Result struct {
	data string
}

func (h *IrisHandler) StartCycleMatch(context *gin.Context) {

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
				ResponseJSON(context, http.StatusOK, res.data)
				return
			}
			global.GVars.UartClient.SendMsg(utils.TrimHexStr(constant.IrisDataMatch))
			go getMatchRes(res)
			time.Sleep(6 * time.Second)
		}
	}
	ResponseJSON(context, http.StatusOK, res)
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
		}
	}

}

func (h *IrisHandler) StopCycleMatch(context *gin.Context) {

	if len(global.GVars.StopCycleChan) == 0 {
		global.GVars.StopCycleChan <- struct{}{}
	}
	ResponseOK(context)
}
