package middleware

import (
	"douyin/pkg/jwt"
	"douyin/pkg/log"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"time"
)

func JWTAuth(j *jwt.JWT, logger *log.Logger) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		code := 0
		msg := ""
		// 获取token, 这里接口文件规定了token时query参数传的，没有放在Header里面
		//token := ctx.GetHeader("Authorization")
		token := ctx.Query("token")
		if token == "" {
			code = -1
			msg = "缺少token"
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"status_code": code,
				"status_msg":  msg,
			})
			ctx.Abort()
			return
		}

		claims, err := j.ParseToken(token)
		if err != nil {
			code = -1
			msg = "解析token时发生了错误"
		} else if time.Now().Unix() > claims.ExpiresAt.Unix() {
			code = -1
			msg = "token已过期"
		}
		if code != 0 {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"status_code": code,
				"status_msg":  msg,
			})
			ctx.Abort()
			return
		}
		ctx.Set("claims", claims)
		fmt.Printf("asdasd")
		recoveryLoggerFunc(ctx, logger)
		ctx.Next()
	}
}

func recoveryLoggerFunc(ctx *gin.Context, logger *log.Logger) {
	userInfo := ctx.MustGet("claims").(*jwt.MyCustomClaims)
	logger.NewContext(ctx, zap.String("UserId", userInfo.UserId))
}
