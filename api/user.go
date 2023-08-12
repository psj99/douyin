package api

import (
	"douyin/service"
	"douyin/service/types/request"
	"douyin/service/types/response"
	"douyin/utils"

	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func UserRegister(ctx *gin.Context) {
	// 绑定JSON到结构体
	req := &request.UserRegisterReq{}
	err := ctx.ShouldBind(req)
	if err != nil {
		utils.ZapLogger.Errorf("ShouldBind err: %v", err)
		ctx.JSON(http.StatusBadRequest, &response.CommonResp{
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
		ctx.JSON(httpCode, &response.CommonResp{
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
		ctx.JSON(http.StatusBadRequest, &response.CommonResp{
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
		ctx.JSON(httpCode, &response.CommonResp{
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

func UserInfo(ctx *gin.Context) {
	// 从请求中读取目标用户ID并与token比对
	target_id := ctx.Query("user_id")
	user_id, ok := ctx.Get("user_id")
	if !ok || target_id != strconv.FormatUint(uint64(user_id.(uint)), 10) {
		utils.ZapLogger.Errorf("UserInfo err: 查询目标与请求用户不同")
		ctx.JSON(http.StatusUnauthorized, &response.CommonResp{
			Status_Code: -1,
			Status_Msg:  "无权读取",
		})
		return
	}

	// 调用获取用户信息
	req := &request.UserInfoReq{User_ID: target_id}
	resp, err := service.UserInfo(ctx, req)
	if err != nil {
		utils.ZapLogger.Errorf("UserInfo err: %v", err)
		ctx.JSON(http.StatusInternalServerError, &response.CommonResp{
			Status_Code: -1,
			Status_Msg:  "获取失败: " + err.Error(),
		})
		return
	}

	// 读取成功
	resp.Status_Code = 0
	resp.Status_Msg = "获取成功"
	ctx.JSON(http.StatusOK, resp)
}
