package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"net/http"
)

func BindJsonAndValid(ctx *gin.Context, bindJson interface{}) error {
	err := ctx.ShouldBindJSON(&bindJson)
	if err != nil {
		return err
	}

	validate := validator.New()
	err = validate.Struct(bindJson)
	if err != nil {
		return err
	}

	return nil
}

func ResponseJSON(c *gin.Context, statusCode int, responseBody interface{}) {
	// 如果 responseBody 是error类型，则返回500
	var _statusCode int
	var _responseBody interface{}
	if value, ok := responseBody.(error); ok {
		_statusCode = http.StatusInternalServerError
		_responseBody = gin.H{
			"msg": value.Error(),
		}
	} else {
		_statusCode = statusCode
		_responseBody = responseBody
	}
	c.JSON(_statusCode, _responseBody)
}

func ResponseOK(c *gin.Context) {
	ResponseJSON(c, http.StatusOK, gin.H{
		"msg": "ok",
	})
}

func ResponseBadRequest(c *gin.Context, msg string) {
	ResponseJSON(c, http.StatusBadRequest, gin.H{
		"msg": msg,
	})
}

func ResponseServerErrorWithData(c *gin.Context, err string, responseBody interface{}) {
	c.JSON(http.StatusInternalServerError, gin.H{
		"msg":  err,
		"data": responseBody,
	})
}

func ResponseServerError(c *gin.Context, err string) {
	c.JSON(http.StatusInternalServerError, gin.H{
		"msg": err,
	})
}

func ResponseUserUnauthorized(c *gin.Context, err string) {
	c.JSON(http.StatusUnauthorized, gin.H{
		"msg": err,
	})
}

func ResponseNoPermission(c *gin.Context) {
	ResponseJSON(c, http.StatusForbidden, gin.H{
		"msg": "当前用户无权执行此操作",
	})
}

func ResponseNoRegister(c *gin.Context) {
	ResponseJSON(c, http.StatusForbidden, gin.H{
		"msg": "当前用户未注册",
	})
}

func ResponseInvalidJSONBody(c *gin.Context) {
	ResponseBadRequest(c, "JSON body 格式错误")
}
