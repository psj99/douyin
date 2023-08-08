package routes

import (
	"douyin/api"
	"github.com/gin-gonic/gin"
)

func NewRouter() *gin.Engine {
	ginRouter := gin.Default()
	app := ginRouter.Group("/douyin")
	{
		app.GET("ping", func(context *gin.Context) {
			context.JSON(200, "success")
		})
	}

	userApi := app.Group("user")
	{
		userApi.POST("/register", api.UserRegister)
		userApi.POST("/login", api.UserLogin)
	}

	return ginRouter
}
