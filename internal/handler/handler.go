package handler

import (
	"douyin/pkg/jwt"
	"douyin/pkg/log"
	"fmt"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	logger *log.Logger
}

func NewHandler(logger *log.Logger) *Handler {
	return &Handler{
		logger: logger,
	}
}

func GetUserIdFromCtx(ctx *gin.Context) string {
	v, exists := ctx.Get("claims")
	if !exists {
		fmt.Println("不存在claims")
		return ""
	}
	return v.(*jwt.MyCustomClaims).UserId
}
