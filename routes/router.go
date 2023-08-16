package routes

import (
	"douyin/api"
	"douyin/utils"

	"net/http"

	"github.com/gin-gonic/gin"
)

func NewRouter() *gin.Engine {
	ginRouter := gin.Default()
	rootAPI := ginRouter.Group("/douyin")
	{
		rootAPI.GET("/ping/", func(context *gin.Context) {
			context.JSON(http.StatusOK, "success")
		})

		rootAPI.GET("/feed/", utils.MiddlewareRateLimit(10, 1), api.GETFeed) // 应用限流中间件 最大10次/秒 每秒恢复1次

		userAPI := rootAPI.Group("user")
		{
			userAPI.POST("/register/", api.POSTUserRegister)
			userAPI.POST("/login/", api.POSTUserLogin)
			userAPI.GET("/", utils.MiddlewareAuth(), api.GETUserInfo) // 应用jwt鉴权中间件
		}

		publishAPI := rootAPI.Group("publish")
		{
			publishAPI.POST("/action/", utils.MiddlewareAuth(), api.POSTPublish) // 应用jwt鉴权中间件
			publishAPI.GET("/list/", utils.MiddlewareAuth(), api.GETPublishList) // 应用jwt鉴权中间件
		}

		favoriteAPI := rootAPI.Group("favorite")
		{
			favoriteAPI.POST("/action/", utils.MiddlewareAuth(), api.POSTFavorite) // 应用jwt鉴权中间件
			favoriteAPI.GET("/list/", utils.MiddlewareAuth(), api.GETFavoriteList) // 应用jwt鉴权中间件
		}

		commentAPI := rootAPI.Group("comment")
		{
			commentAPI.POST("/action/", utils.MiddlewareAuth(), api.POSTComment) // 应用jwt鉴权中间件
			commentAPI.GET("/list/", utils.MiddlewareAuth(), api.GETCommentList) // 应用jwt鉴权中间件
		}

		relationAPI := rootAPI.Group("relation")
		{
			relationAPI.POST("/action/", utils.MiddlewareAuth(), api.POSTFollow)            // 应用jwt鉴权中间件
			relationAPI.GET("/follow/list/", utils.MiddlewareAuth(), api.GETFollowList)     // 应用jwt鉴权中间件
			relationAPI.GET("/follower/list/", utils.MiddlewareAuth(), api.GETFollowerList) // 应用jwt鉴权中间件
			relationAPI.GET("/friend/list/", utils.MiddlewareAuth(), api.GETFriendList)     // 应用jwt鉴权中间件
		}

		messageAPI := rootAPI.Group("message")
		{
			messageAPI.POST("/action/", utils.MiddlewareAuth(), api.POSTMessage) // 应用jwt鉴权中间件
			messageAPI.GET("/chat/", utils.MiddlewareAuth(), api.GETMessageList) // 应用jwt鉴权中间件
		}
	}
	return ginRouter
}
