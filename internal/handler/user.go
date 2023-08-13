package handler

import (
	"douyin/internal/pkg/request"
	"douyin/internal/pkg/resp"
	"douyin/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UserHandler interface {
	Register(ctx *gin.Context)
	Login(ctx *gin.Context)
	GetUserInfo(ctx *gin.Context)
}

type userHandler struct {
	*Handler
	userService service.UserService
}

func (u userHandler) Register(ctx *gin.Context) {
	// 绑定JSON到结构体
	req := &request.UserRegisterReq{}
	err := ctx.ShouldBind(req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, &resp.UserRegisterResp{
			Response: resp.ResponseErr("输入有误, 请重新输入"),
		})
		return
	}
	userId, token, err := u.userService.Register(ctx, req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, &resp.UserRegisterResp{
			Response: resp.ResponseErr("注册失败, " + err.Error()),
		})
		return
	}
	// 注册成功
	ctx.JSON(http.StatusOK, &resp.UserRegisterResp{
		Response: resp.ResponseOK(),
		UserId:   int64(userId),
		Token:    token,
	})
}

func (u userHandler) Login(ctx *gin.Context) {
	// 绑定JSON到结构体
	req := &request.UserLoginReq{}
	err := ctx.ShouldBind(req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, &resp.UserRegisterResp{
			Response: resp.ResponseErr("输入有误, 请重新输入"),
		})
		return
	}

	// 调用用户登录处理
	userId, token, err := u.userService.Login(ctx, req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, &resp.UserRegisterResp{
			Response: resp.ResponseErr("登录失败, " + err.Error()),
		})
		return
	}
	// 登录成功
	ctx.JSON(http.StatusOK, &resp.UserRegisterResp{
		Response: resp.ResponseOK(),
		UserId:   int64(userId),
		Token:    token,
	})
}

func (u userHandler) GetUserInfo(ctx *gin.Context) {
	//TODO implement me
	userId := GetUserIdFromCtx(ctx)
	u.logger.Info("userId:" + userId)
	if userId == "" {
		ctx.JSON(http.StatusUnauthorized, resp.UserInfoResp{
			Response: resp.ResponseErr("请先登录"),
			UserInfo: nil,
		})
		return
	}
	id, _ := strconv.Atoi(userId)
	userInfo, err := u.userService.GetUserInfo(ctx, uint(id))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, resp.UserInfoResp{
			Response: resp.ResponseErr("服务器出错了"),
			UserInfo: nil,
		})
		return
	}

	ctx.JSON(http.StatusOK, resp.UserInfoResp{
		Response: resp.ResponseOK(),
		UserInfo: userInfo,
	})
}

func NewUserHandler(handler *Handler, userService service.UserService) UserHandler {
	return &userHandler{
		Handler:     handler,
		userService: userService,
	}
}
