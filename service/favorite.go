package service

import (
	"douyin/repository/dao"
	"douyin/service/types/request"
	"douyin/service/types/response"
	"douyin/utils"
	"douyin/utils/oss"

	"context"
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
)

func Favorite(ctx *gin.Context, req *request.FavoriteReq) (resp *response.FavoriteResp, err error) {
	// 获取请求用户ID
	Me_ID, ok := ctx.Get("user_id")
	if !ok {
		utils.ZapLogger.Errorf("ctx.Get (user_id) err: inaccessible")
		return nil, errors.New("无法获取user_id")
	}

	// 存储点赞信息
	action_type, err := strconv.ParseUint(req.Action_Type, 10, 64)
	if err != nil {
		utils.ZapLogger.Errorf("ParseUint err: %v", err)
		return nil, err
	}
	video_id, err := strconv.ParseUint(req.Video_ID, 10, 64)
	if err != nil {
		utils.ZapLogger.Errorf("ParseUint err: %v", err)
		return nil, err
	}

	if action_type == 1 {
		err = dao.CreateFavorite(context.TODO(), Me_ID.(uint), uint(video_id))
		if err != nil {
			utils.ZapLogger.Errorf("CreateFavorite err: %v", err)
			return nil, err
		}
	} else if action_type == 2 {
		err = dao.DeleteFavorite(context.TODO(), Me_ID.(uint), uint(video_id))
		if err != nil {
			utils.ZapLogger.Errorf("DeleteFavorite err: %v", err)
			return nil, err
		}
	} else {
		utils.ZapLogger.Errorf("Invalid action_type err: %v", action_type)
		return nil, errors.New("操作类型有误")
	}

	return &response.FavoriteResp{}, nil
}

func FavoriteList(ctx *gin.Context, req *request.FavoriteListReq) (resp *response.FavoriteListResp, err error) {
	// 获取请求用户ID
	Me_ID, ok := ctx.Get("user_id")
	if !ok {
		utils.ZapLogger.Errorf("ctx.Get (user_id) err: inaccessible")
		return nil, errors.New("无法获取user_id")
	}

	// 获取目标用户信息
	userID, err := strconv.ParseUint(req.User_ID, 10, 64) // string转十进制uint64
	if err != nil {
		utils.ZapLogger.Errorf("ParseUint err: %v", err)
		return nil, err
	}
	user, err := dao.FindUserByID(context.TODO(), uint(userID))
	if err != nil {
		utils.ZapLogger.Errorf("FindUserByID err: %v", err)
		return nil, err
	}

	// 读取视频列表
	// 临时方案 亟待优化 //TODO
	resp = &response.FavoriteListResp{}
	for _, video := range user.Favorites {
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
		isFavorite, err = dao.CheckFavorite(context.TODO(), Me_ID.(uint), video.ID)
		if err != nil {
			isFavorite = false
			utils.ZapLogger.Errorf("CheckFavorite err: %v", err)
		}
		videoInfo.Is_Favorite = isFavorite

		// 获取作者信息
		author, err := dao.FindUserByID(context.TODO(), video.UserID)
		if err != nil {
			utils.ZapLogger.Errorf("FindUserByID (author) err: %v", err)
			continue // 跳过此条视频
		}
		followCount := uint(len(author.Follows))     // 统计关注数
		followerCount := uint(len(author.Followers)) // 统计粉丝数
		workCount := uint(len(author.Works))         // 统计作品数
		favoriteCount := uint(len(author.Favorites)) // 统计点赞数

		// 统计获赞数
		var favoritedCount uint = 0
		for _, v := range author.Works {
			favoritedCount += uint(len(v.Favorited))
		}

		// 是否关注
		isFollow, err := dao.CheckFollow(context.TODO(), Me_ID.(uint), uint(author.ID))
		if err != nil {
			isFollow = false
			utils.ZapLogger.Errorf("CheckFollow err: %v", err)
		}

		// 视频信息中加入作者信息
		videoInfo.Author = response.User{
			ID:               author.ID,
			Name:             author.Username,
			Follow_Count:     followCount,
			Follower_Count:   followerCount,
			Is_Follow:        isFollow,
			Avatar:           author.Avatar,
			Background_Image: author.BackgroundImage,
			Signature:        author.Signature,
			Total_Favorited:  strconv.FormatUint(uint64(favoritedCount), 10),
			Work_Count:       workCount,
			Favorite_Count:   favoriteCount,
		}

		// 将该视频加入列表
		resp.Video_List = append(resp.Video_List, videoInfo)
	}

	return resp, nil
}
