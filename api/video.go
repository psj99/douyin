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

func GETVideoFeed(ctx *gin.Context) {
	// 获取可选参数
	latest_time := ctx.Query("latest_time")
	if latest_time == "" || latest_time == "0" {
		latest_time = strconv.FormatInt(time.Now().Unix(), 10) // 空值表示当前时间
	} else {
		latest_time = latest_time[:len(latest_time)-3] // API文档有误 请求实为毫秒时间戳 故在此转换
	}

	// 调用获取用户信息
	req := &request.FeedReq{Latest_Time: latest_time, Token: ctx.Query("token")}
	resp, err := service.VideoFeed(ctx, req)
	if err != nil {
		utils.ZapLogger.Errorf("VideoFeed err: %v", err)
		ctx.JSON(http.StatusInternalServerError, &response.Status{
			Status_Code: -1,
			Status_Msg:  "获取失败: " + err.Error(),
		})
		return
	}

	// 读取成功
	status := response.Status{Status_Code: 0, Status_Msg: "获取成功"}
	resp.Status = status
	ctx.JSON(http.StatusOK, resp)
}

func POSTVideoPublish(ctx *gin.Context) {
	// 绑定JSON到结构体
	req := &request.PublishReq{}
	err := ctx.ShouldBind(req)
	if err != nil {
		utils.ZapLogger.Errorf("ShouldBind err: %v", err)
		ctx.JSON(http.StatusBadRequest, &response.Status{
			Status_Code: -1,
			Status_Msg:  "发布失败: " + err.Error(),
		})
		return
	}

	// 获取上传的文件
	file, err := ctx.FormFile("data") // 只获取请求中附带的第一个文件
	if err != nil {
		utils.ZapLogger.Errorf("FormFile err: %v", err)
		ctx.JSON(http.StatusBadRequest, &response.Status{
			Status_Code: -1,
			Status_Msg:  "发布失败: " + err.Error(),
		})
		return
	}

	// 调用投稿处理
	resp, err := service.VideoPublish(ctx, req, file)
	if err != nil {
		utils.ZapLogger.Errorf("VideoPublish err: %v", err)
		ctx.JSON(http.StatusInternalServerError, &response.Status{
			Status_Code: -1,
			Status_Msg:  "发布失败: " + err.Error(),
		})
		return
	}

	// 发布成功
	status := response.Status{Status_Code: 0, Status_Msg: "发布成功"}
	resp.Status = status
	ctx.JSON(http.StatusOK, resp)
}

func GETVideoPublishList(ctx *gin.Context) {
	// 绑定JSON到结构体
	req := &request.PublishListReq{}
	err := ctx.ShouldBind(req)
	if err != nil {
		utils.ZapLogger.Errorf("ShouldBind err: %v", err)
		ctx.JSON(http.StatusBadRequest, &response.Status{
			Status_Code: -1,
			Status_Msg:  "获取失败: " + err.Error(),
		})
		return
	}

	/*// 从请求中读取目标用户ID并与token比对
	user_id, ok := ctx.Get("user_id")
	if !ok || req.User_ID != strconv.FormatUint(uint64(user_id.(uint)), 10) {
		utils.ZapLogger.Errorf("GETVideoPublishList err: 查询目标与请求用户不同")
		ctx.JSON(http.StatusUnauthorized, &response.Status{
			Status_Code: -1,
			Status_Msg:  "无权获取",
		})
		return
	}*/

	// 调用获取用户信息
	resp, err := service.VideoPublishList(ctx, req)
	if err != nil {
		utils.ZapLogger.Errorf("VideoPublishList err: %v", err)
		ctx.JSON(http.StatusInternalServerError, &response.Status{
			Status_Code: -1,
			Status_Msg:  "获取失败: " + err.Error(),
		})
		return
	}

	// 读取成功
	status := response.Status{Status_Code: 0, Status_Msg: "获取成功"}
	resp.Status = status
	ctx.JSON(http.StatusOK, resp)
}
