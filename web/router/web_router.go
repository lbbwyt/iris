package router

import (
	"github.com/gin-gonic/gin"
	"iris/web/handler"
)

func InitRouter() (*gin.Engine, error) {
	r := gin.Default()
	gin.SetMode("debug")

	r.GET("/ping", func(context *gin.Context) {
		context.JSON(200, gin.H{
			"msg": "pong",
		})
	})

	api := r.Group("/iris")
	//虹膜是被
	_ = RegisterIrisRouter(api)
	return r, nil
}

func RegisterIrisRouter(rg *gin.RouterGroup) error {

	handler := handler.NewIrisHandler()

	// routes
	{
		rg.POST("/actions/register", handler.RegisterIris) // 虹膜注册
		rg.POST("/actions/match", handler.MatchIris)       // 虹膜识别
		rg.DELETE("/delete", handler.DeleteAllUser)        // 删除所有用户
		rg.POST("/mode", handler.ChangeMode)               // 切换虹膜预览模式
		rg.POST("/cycle/start", handler.StartCycleMatch)   // 开始循环识别
		rg.POST("/cycle/stop", handler.StopCycleMatch)     // 停止循环识别
	}
	return nil
}
