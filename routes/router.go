package routes

import (
	"douyin/api"
	"douyin/utils"

	"net/http"

	"github.com/gin-gonic/gin"
)

func NewRouter() *gin.Engine {
	ginRouter := gin.Default()
	rootApi := ginRouter.Group("/douyin")
	{
		rootApi.GET("/ping/", func(context *gin.Context) {
			context.JSON(http.StatusOK, "success")
		})

		rootApi.GET("/feed/", utils.MiddlewareRateLimit(10, 1), api.VideoFeed) // 应用限流中间件 最大10次/秒 每秒恢复1次

		userApi := rootApi.Group("user")
		{
			userApi.POST("/register/", api.UserRegister)
			userApi.POST("/login/", api.UserLogin)
		}
	}
	return ginRouter
}
