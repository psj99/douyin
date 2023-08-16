package service

import (
	"douyin/repository/dao"
	"douyin/repository/model"
	"douyin/service/types/request"
	"douyin/service/types/response"
	"douyin/utils"
	"douyin/utils/oss"

	"context"
	"errors"
	"mime/multipart"
	"strconv"

	"github.com/gin-gonic/gin"
)

// 发布视频
func Publish(ctx *gin.Context, req *request.PublishReq, file *multipart.FileHeader) (resp *response.PublishResp, err error) {
	// 首先尝试打开文件
	videoStream, err := file.Open()
	if err != nil {
		utils.ZapLogger.Errorf("file.Open err: %v", err)
		return nil, err
	}
	defer videoStream.Close() // 不保证自动关闭成功

	// 准备要存储的内容
	userID, ok := ctx.Get("user_id")
	if !ok {
		utils.ZapLogger.Errorf("ctx.Get (user_id) err: %v", err)
		return nil, errors.New("无法获取user_id")
	}
	video := &model.Video{
		Title:  req.Title,
		UserID: userID.(uint),
	}

	// 存储视频信息
	video, err = dao.CreateVideo(context.TODO(), video)
	if err != nil {
		utils.ZapLogger.Errorf("CreateVideo err: %v", err)
		return nil, err
	}

	// 上传视频数据(封面为默认)
	videoID := strconv.FormatUint(uint64(video.ID), 10)
	err = oss.UploadVideoStream(context.TODO(), videoID, videoStream, file.Size)
	if err != nil {
		utils.ZapLogger.Errorf("UploadVideoStream err: %v", err)
		return nil, err
	}

	// 创建更新封面异步任务
	go func() {
		oss.UpdateCover(context.TODO(), videoID) // 不保证自动更新成功
	}()

	return &response.PublishResp{}, nil
}

// 获取发布列表
func PublishList(ctx *gin.Context, req *request.PublishListReq) (resp *response.PublishListResp, err error) {
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

	// 组装作者(目标用户)信息
	// 临时方案 亟待优化 //TODO
	followCount := uint(len(user.Follows))     // 统计关注数
	followerCount := uint(len(user.Followers)) // 统计粉丝数
	workCount := uint(len(user.Works))         // 统计作品数
	favoriteCount := uint(len(user.Favorites)) // 统计点赞数

	// 统计获赞数
	var favoritedCount uint = 0
	for _, video := range user.Works {
		favoritedCount += uint(len(video.Favorited))
	}

	// 是否关注
	isFollow := dao.CheckFollow(context.TODO(), Me_ID.(uint), uint(userID))

	// 作者信息
	userInfo := response.User{
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
	}

	// 读取视频列表
	// 临时方案 亟待优化 //TODO
	resp = &response.PublishListResp{}
	for _, video := range user.Works {
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
		videoInfo.Is_Favorite = dao.CheckFavorite(context.TODO(), Me_ID.(uint), video.ID)

		// 视频信息中加入作者信息
		videoInfo.Author = userInfo

		// 将该视频加入列表
		resp.Video_List = append(resp.Video_List, videoInfo)
	}

	return resp, nil
}
