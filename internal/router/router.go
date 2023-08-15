package router

import (
	"douyin/internal/handler"
	"douyin/internal/pkg/middleware"
	"douyin/pkg/jwt"
	"douyin/pkg/log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func NewRouter(
	logger *log.Logger,
	jwt *jwt.JWT,
	userHandler handler.UserHandler,
	videoHandler handler.VideoHandler,
) *gin.Engine {
	gin.SetMode(gin.DebugMode)
	r := gin.Default()

	//r.Use(
	//middleware.CORSMiddleware(),
	//middleware.ResponseLogMiddleware(logger),
	//middleware.RequestLogMiddleware(logger),
	//middleware.SignMiddleware(log),
	//)

	// 不需要登陆
	app := r.Group("/douyin")
	{
		app.GET("/", func(ctx *gin.Context) {
			logger.WithContext(ctx).Info("hello")
			ctx.JSON(http.StatusOK, map[string]interface{}{
				"say": "Hi DouYin!",
			})
		})

		app.POST("/user/register", userHandler.Register)
		app.POST("/user/login", userHandler.Login)
		app.GET("/feed", videoHandler.Feed)
	}

	// 需要登录
	auth := app.Group("/")
	auth.Use(middleware.JWTAuth(jwt, logger))
	auth.GET("/user", userHandler.GetUserInfo)
	auth.POST("/publish/action", videoHandler.PublishAction)
	auth.GET("/publish/list", videoHandler.PublishList)
	return r
}
