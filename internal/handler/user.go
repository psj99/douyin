package handler

import (
	"douyin/internal/pkg/request"
	"douyin/internal/pkg/response"
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
		ctx.JSON(http.StatusBadRequest, &response.UserRegisterResp{
			StatusCode: 7,
			StatusMsg:  "输入有误",
		})
		return
	}

	resp, err := u.userService.Register(ctx, req)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, &response.UserRegisterResp{
			StatusCode: 7,
			StatusMsg:  "注册失败",
		})
		return
	}

	// 注册成功
	ctx.JSON(http.StatusOK, resp)
}

func (u userHandler) Login(ctx *gin.Context) {
	// 绑定JSON到结构体
	req := &request.UserLoginReq{}
	err := ctx.ShouldBind(req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, &response.UserRegisterResp{
			StatusCode: 7,
			StatusMsg:  "输入有误",
		})
		return
	}

	// 调用用户登录处理
	resp, err := u.userService.Login(ctx, req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, &response.UserLoginResp{
			StatusCode: -1,
			StatusMsg:  "登录失败: " + err.Error(),
		})
		return
	}

	// 登录成功
	ctx.JSON(http.StatusOK, resp)
}

func (u userHandler) GetUserInfo(ctx *gin.Context) {
	//TODO implement me
	userId := GetUserIdFromCtx(ctx)
	u.logger.Info("userId:" + userId)
	if userId == "" {
		ctx.JSON(http.StatusUnauthorized, response.UserInfoResp{
			StatusCode: 7,
			StatusMsg:  "请先登录",
			User:       nil,
		})
		return
	}
	id, _ := strconv.Atoi(userId)
	userInfo, err := u.userService.GetUserInfo(ctx, uint(id))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.UserInfoResp{
			StatusCode: 7,
			StatusMsg:  "服务器出错了",
			User:       nil,
		})
		return
	}

	ctx.JSON(http.StatusOK, response.UserInfoResp{
		StatusCode: 0,
		StatusMsg:  "查询成功",
		User:       userInfo,
	})
}

func NewUserHandler(handler *Handler, userService service.UserService) UserHandler {
	return &userHandler{
		Handler:     handler,
		userService: userService,
	}
}
