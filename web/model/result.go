package model

import "net/http"

type Result struct {
	Success bool        `json:"result"`
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func ResultOk(data interface{}) Result {
	return Result{
		Success: true,
		Code:    http.StatusOK,
		Message: "success",
		Data:    data,
	}
}

func ResultError(errorCode int, msg string) Result {
	return Result{
		Success: false,
		Code:    errorCode,
		Message: msg,
		Data:    nil,
	}
}
