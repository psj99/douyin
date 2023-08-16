package service

import (
	"douyin/repository/dao"
	"douyin/service/types/response"
	"douyin/utils"

	"context"
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
)

// 读取指定用户信息 返回用户信息响应结构体
func readUserInfo(ctx *gin.Context, user_id uint) (userInfo *response.User, err error) {
	// 获取请求用户ID
	Me_ID, ok := ctx.Get("user_id")
	if !ok {
		utils.ZapLogger.Errorf("ctx.Get (user_id) err: inaccessible")
		return nil, errors.New("无法获取user_id")
	}

	// 读取目标用户信息
	user, err := dao.FindUserByID(context.TODO(), user_id)
	if err != nil {
		utils.ZapLogger.Errorf("FindUserByID err: %v", err)
		return nil, err
	}
	followCount := uint(len(user.Follows))     // 统计关注数
	followerCount := uint(len(user.Followers)) // 统计粉丝数
	workCount := uint(len(user.Works))         // 统计作品数
	favoriteCount := uint(len(user.Favorites)) // 统计点赞数

	// 统计获赞数
	var favoritedCount uint = 0
	for _, v := range user.Works {
		favoritedCount += uint(len(v.Favorited))
	}

	// 是否关注
	isFollow := false
	if Me_ID != "" {
		isFollow = dao.CheckFollow(context.TODO(), Me_ID.(uint), uint(user.ID))
	}

	return &response.User{
		ID:               user.ID,
		Name:             user.Username,
		Follow_Count:     followCount,
		Follower_Count:   followerCount,
		Is_Follow:        isFollow,
		Avatar:           user.Avatar,
		Background_Image: user.BackgroundImage,
		Signature:        user.Signature,
		Total_Favorited:  strconv.FormatUint(uint64(favoritedCount), 10),
		Work_Count:       workCount,
		Favorite_Count:   favoriteCount,
	}, nil
}
