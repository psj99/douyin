package handler

import (
	"douyin/internal/service"

	"github.com/gin-gonic/gin"
)

type UserHandler interface {
	Register(ctx *gin.Context)
	Login(ctx *gin.Context)
	GetProfile(ctx *gin.Context)
	UpdateProfile(ctx *gin.Context)
}

type userHandler struct {
	*Handler
	userService service.UserService
}

func NewUserHandler(handler *Handler, userService service.UserService) UserHandler {
	return &userHandler{
		Handler:     handler,
		userService: userService,
	}
}

func (h *userHandler) Register(ctx *gin.Context) {
	//TODO: implement me
	panic("implement me")
}

func (h *userHandler) Login(ctx *gin.Context) {
	//TODO: implement me
	panic("implement me")
}

func (h *userHandler) GetProfile(ctx *gin.Context) {
	//TODO: implement me
	panic("implement me")
}

func (h *userHandler) UpdateProfile(ctx *gin.Context) {
	//TODO: implement me
	panic("implement me")
}
