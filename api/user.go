package api

import (
	"douyin/service"
	"douyin/service/types/request"
	"douyin/service/types/response"
	"douyin/utils"

	"net/http"

	"github.com/gin-gonic/gin"
)

func UserRegister(ctx *gin.Context) {
	// 绑定JSON到结构体
	req := &request.UserRegisterReq{}
	err := ctx.ShouldBind(req)
	if err != nil {
		utils.ZapLogger.Errorf("ShouldBind err: %v", err)
		ctx.JSON(http.StatusBadRequest, &response.UserRegisterResp{
			Status_Code: -1,
			Status_Msg:  "注册失败: " + err.Error(),
		})
		return
	}

	// 调用用户注册处理
	resp, err := service.UserRegister(ctx, req)
	if err != nil {
		utils.ZapLogger.Errorf("UserRegister err: %v", err)
		var httpCode int
		if err == service.ErrorUserExists {
			httpCode = http.StatusConflict
		} else {
			httpCode = http.StatusInternalServerError
		}
		ctx.JSON(httpCode, &response.UserRegisterResp{
			Status_Code: -1,
			Status_Msg:  "注册失败: " + err.Error(),
		})
		return
	}

	// 注册成功
	resp.Status_Code = 0
	resp.Status_Msg = "注册成功"
	ctx.JSON(http.StatusOK, resp)
}

func UserLogin(ctx *gin.Context) {
	// 绑定JSON到结构体
	req := &request.UserLoginReq{}
	err := ctx.ShouldBind(req)
	if err != nil {
		utils.ZapLogger.Errorf("ShouldBind err: %v", err)
		ctx.JSON(http.StatusBadRequest, &response.UserLoginResp{
			Status_Code: -1,
			Status_Msg:  "登录失败: " + err.Error(),
		})
		return
	}

	// 调用用户登录处理
	resp, err := service.UserLogin(ctx, req)
	if err != nil {
		utils.ZapLogger.Errorf("UserLogin err: %v", err)
		var httpCode int
		if err == service.ErrorUserNotExists || err == service.ErrorWrongPassword {
			httpCode = http.StatusUnauthorized
		} else {
			httpCode = http.StatusInternalServerError
		}
		ctx.JSON(httpCode, &response.UserLoginResp{
			Status_Code: -1,
			Status_Msg:  "登录失败: " + err.Error(),
		})
		return
	}

	// 登录成功
	resp.Status_Code = 0
	resp.Status_Msg = "登录成功"
	ctx.JSON(http.StatusOK, resp)
}
