package api

import (
	"douyin/pkg/types/request"
	"douyin/pkg/types/response"
	"douyin/pkg/utils"
	"douyin/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

func UserRegister(ctx *gin.Context) {
	var req request.UserRegisterReq
	err := ctx.ShouldBind(&req)
	// 参数校验
	if err != nil {
		utils.ZapLogger.Infoln(err)
		return
	}
	userSrv := service.GetUserSrv()
	resp, err := userSrv.UserRegister(ctx, &req)
	if err != nil {
		utils.ZapLogger.Infof("err: %v", err)
		ctx.JSON(http.StatusBadRequest, &response.UserLoginResp{
			Code: -1,
			Msg:  "注册失败: " + err.Error(),
		})
	}
	ctx.JSON(http.StatusOK, resp)
}

func UserLogin(ctx *gin.Context) {
	var req request.UserLoginReq
	err := ctx.ShouldBind(&req)
	// 参数校验
	if err != nil {
		utils.ZapLogger.Infoln(err)
		ctx.JSON(http.StatusBadRequest, &response.UserLoginResp{
			Code: -1,
			Msg:  "请求参数错误",
		})
		return
	}
	userSrv := service.GetUserSrv()
	resp, err := userSrv.UserLogin(ctx, &req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, &response.UserLoginResp{
			Code: -1,
			Msg:  "登录失败: " + err.Error(),
		})
	}
	ctx.JSON(http.StatusOK, resp)

}
