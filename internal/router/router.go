package router

import (
	"douyin/internal/handler"
	"douyin/pkg/jwt"
	"douyin/pkg/log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func NewRouter(
	logger *log.Logger,
	jwt *jwt.JWT,
	userHandler handler.UserHandler,
) *gin.Engine {
	// gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	r.Use(
	// middleware.CORSMiddleware(),
	// middleware.ResponseLogMiddleware(logger),
	// middleware.RequestLogMiddleware(logger),
	//middleware.SignMiddleware(log),
	)

	// No route group has permission
	noAuthRouter := r.Group("/")
	{

		noAuthRouter.GET("/", func(ctx *gin.Context) {
			logger.WithContext(ctx).Info("hello")
			ctx.JSON(http.StatusOK, map[string]interface{}{
				"say": "Hi Nunu!",
			})
		})

		noAuthRouter.POST("/register", userHandler.Register)
		noAuthRouter.POST("/login", userHandler.Login)
	}

	// // Non-strict permission routing group
	// noStrictAuthRouter := r.Group("/").Use(middleware.NoStrictAuth(jwt, logger))
	// {
	// 	noStrictAuthRouter.GET("/user", userHandler.GetProfile)
	// }

	// // Strict permission routing group
	// strictAuthRouter := r.Group("/").Use(middleware.StrictAuth(jwt, logger))
	// {
	// 	strictAuthRouter.PUT("/user", userHandler.UpdateProfile)
	// }

	return r

}
