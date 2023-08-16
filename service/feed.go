package service

import (
	"douyin/repository/dao"
	"douyin/service/types/request"
	"douyin/service/types/response"
	"douyin/utils"
	"douyin/utils/oss"

	"context"
	"strconv"

	"github.com/gin-gonic/gin"
)

// 视频流
func Feed(ctx *gin.Context, req *request.FeedReq) (resp *response.FeedResp, err error) {
	// 获取请求用户ID 默认为0(不存在)
	var Me_ID uint = 0
	if req.Token != "" {
		// 解析/校验token (自动验证有效期等)
		claims, err := utils.ParseToken(req.Token)
		if err == nil { // 若成功登录
			// 提取user_id和username
			Me_ID = uint(claims["user_id"].(float64)) // token中解析数字默认float64
		}
	}

	// 读取视频列表
	latest_time, err := strconv.ParseInt(req.Latest_Time, 10, 64)
	if err != nil {
		utils.ZapLogger.Errorf("ParseInt err: %v", err)
		return nil, err
	}
	videos, err := dao.FindVideosByUpdatedAt(context.TODO(), latest_time, false, 30) // 倒序向过去查找 最多30条
	if err != nil {
		utils.ZapLogger.Errorf("FindVideosByUpdatedAt err: %v", err)
		return nil, err
	}

	// 临时方案 亟待优化 //TODO
	resp = &response.FeedResp{
		Next_Time: 0, // 本次返回的视频中发布最早的时间 默认为无效
	}
	if len(videos) > 0 { // 如果查找结果中有视频
		resp.Next_Time = uint(videos[len(videos)-1].UpdatedAt.Unix() * 1000) // 更新该时间戳 API文档有误 请求实为毫秒时间戳 故在此转换
	}

	// 向列表中添加视频
	for _, video := range videos {
		// 读取视频信息
		videoInfo := response.Video{
			ID:             video.ID,
			Favorite_Count: uint(len(video.Favorited)),
			Comment_Count:  uint(len(video.Comments)),
			Title:          video.Title,
		}

		// 获取视频及封面URL
		objectID := strconv.FormatUint(uint64(video.ID), 10)
		videoURL, coverURL, err := oss.GetVideo(context.TODO(), objectID)
		if err != nil {
			utils.ZapLogger.Errorf("GetVideo err: %v", err)
			continue // 跳过此条视频
		}
		videoInfo.Play_URL = videoURL
		videoInfo.Cover_URL = coverURL

		// 检查是否点赞
		isFavorite := false
		if Me_ID != 0 {
			isFavorite = dao.CheckFavorite(context.TODO(), Me_ID, video.ID)
		}
		videoInfo.Is_Favorite = isFavorite

		// 读取作者
		user, err := dao.FindUserByID(context.TODO(), video.UserID)
		if err != nil {
			utils.ZapLogger.Errorf("FindUserByID err: %v", err)
			continue // 跳过此条视频
		}
		userInfo := response.User{
			ID:               user.ID,
			Name:             user.Username,
			Follow_Count:     uint(len(user.Follows)),
			Follower_Count:   uint(len(user.Followers)),
			Avatar:           user.Avatar,
			Background_Image: user.BackgroundImage,
			Signature:        user.Signature,
			Work_Count:       uint(len(user.Works)),
			Favorite_Count:   uint(len(user.Favorites)),
		}

		var favoritedCount uint = 0 // 统计获赞数
		for _, video := range user.Works {
			favoritedCount += uint(len(video.Favorited))
		}
		userInfo.Total_Favorited = strconv.FormatUint(uint64(favoritedCount), 10)

		isFollow := false
		if Me_ID != 0 { // 若登录则检查是否关注
			isFollow = dao.CheckFollow(context.TODO(), Me_ID, user.ID)
		}
		userInfo.Is_Follow = isFollow

		// 视频信息中加入作者信息
		videoInfo.Author = userInfo

		// 将该视频加入列表
		resp.Video_List = append(resp.Video_List, videoInfo)
	}

	return resp, nil
}
