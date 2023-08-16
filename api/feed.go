package api

import (
	"douyin/service"
	"douyin/service/types/request"
	"douyin/service/types/response"
	"douyin/utils"

	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func GETFeed(ctx *gin.Context) {
	// 获取可选参数
	latest_time := ctx.Query("latest_time")
	if latest_time == "" || latest_time == "0" {
		latest_time = strconv.FormatInt(time.Now().Unix(), 10) // 空值表示当前时间
	} else {
		latest_time = latest_time[:len(latest_time)-3] // API文档有误 请求实为毫秒时间戳 故在此转换
	}

	// 调用获取用户信息
	req := &request.FeedReq{Latest_Time: latest_time, Token: ctx.Query("token")}
	resp, err := service.Feed(ctx, req)
	if err != nil {
		utils.ZapLogger.Errorf("Feed err: %v", err)
		ctx.JSON(http.StatusInternalServerError, &response.Status{
			Status_Code: -1,
			Status_Msg:  "获取失败: " + err.Error(),
		})
		return
	}

	// 获取成功
	status := response.Status{Status_Code: 0, Status_Msg: "获取成功"}
	resp.Status = status
	ctx.JSON(http.StatusOK, resp)
}
