package service

import (
	"douyin/repository/dao"
	"douyin/service/types/request"
	"douyin/service/types/response"
	"douyin/utils"

	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
)

// 评论/删除评论
func Comment(ctx *gin.Context, req *request.CommentReq) (resp *response.CommentResp, err error) {
	// 获取请求用户ID
	Me_ID, ok := ctx.Get("user_id")
	if !ok {
		utils.ZapLogger.Errorf("ctx.Get (user_id) err: inaccessible")
		return nil, errors.New("无法获取user_id")
	}

	// 存储评论信息
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
		// 创建评论
		comment, err := dao.CreateComment(context.TODO(), Me_ID.(uint), uint(video_id), req.Comment_Text)
		if err != nil {
			utils.ZapLogger.Errorf("CreateComment err: %v", err)
			return nil, err
		}

		// 组装评论信息
		commentInfo := response.Comment{ID: comment.ID, Content: comment.Content}

		// 评论发布时间
		create_month := comment.CreatedAt.Month()
		create_day := comment.CreatedAt.Day()
		commentInfo.Create_Date = fmt.Sprintf("%02d-%02d", create_month, create_day) // mm-dd

		// 临时方案 亟待优化 //TODO
		// 评论作者信息
		author, err := dao.FindUserByID(context.TODO(), Me_ID.(uint))
		if err == nil {
			followCount := uint(len(author.Follows))     // 统计关注数
			followerCount := uint(len(author.Followers)) // 统计粉丝数
			workCount := uint(len(author.Works))         // 统计作品数
			favoriteCount := uint(len(author.Favorites)) // 统计点赞数

			// 统计获赞数
			var favoritedCount uint = 0
			for _, video := range author.Works {
				favoritedCount += uint(len(video.Favorited))
			}

			// 是否关注
			isFollow := dao.CheckFollow(context.TODO(), Me_ID.(uint), Me_ID.(uint))

			userInfo := response.User{
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

			// 评论信息中加入作者信息
			commentInfo.User = userInfo
		} else { // 若查找用户失败则作者为空 不阻止运行
			utils.ZapLogger.Errorf("FindUserByID err: %v", err)
		}

		return &response.CommentResp{
			Comment: commentInfo,
		}, nil
	} else if action_type == 2 {
		comment_id, err := strconv.ParseUint(req.Comment_ID, 10, 64)
		if err != nil {
			utils.ZapLogger.Errorf("ParseUint err: %v", err)
			return nil, err
		}

		err = dao.DeleteComment(context.TODO(), uint(comment_id), true) // 永久删除
		if err != nil {
			utils.ZapLogger.Errorf("DeleteComment err: %v", err)
			return nil, err
		}

		return &response.CommentResp{}, nil // 不返回被删除评论内容
	} else {
		utils.ZapLogger.Errorf("Invalid action_type err: %v", action_type)
		return nil, errors.New("操作类型有误")
	}
}

// 获取评论列表
func CommentList(ctx *gin.Context, req *request.CommentListReq) (resp *response.CommentListResp, err error) {
	// 获取请求用户ID
	Me_ID, ok := ctx.Get("user_id")
	if !ok {
		utils.ZapLogger.Errorf("ctx.Get (user_id) err: inaccessible")
		return nil, errors.New("无法获取user_id")
	}

	// 读取评论列表
	video_id, err := strconv.ParseUint(req.Video_ID, 10, 64)
	if err != nil {
		utils.ZapLogger.Errorf("ParseUint err: %v", err)
		return nil, err
	}
	comments, err := dao.FindCommentsByCreatedAt(context.TODO(), uint(video_id), false)
	if err != nil {
		utils.ZapLogger.Errorf("FindCommentsByCreatedAt err: %v", err)
		return nil, err
	}

	// 临时方案 亟待优化 //TODO
	// 向列表中添加评论
	resp = &response.CommentListResp{}
	for _, comment := range comments {
		// 读取评论信息
		commentInfo := response.Comment{ID: comment.ID, Content: comment.Content}

		// 评论发布时间
		create_month := comment.CreatedAt.Month()
		create_day := comment.CreatedAt.Day()
		commentInfo.Create_Date = fmt.Sprintf("%02d-%02d", create_month, create_day) // mm-dd

		// 评论作者信息
		author, err := dao.FindUserByID(context.TODO(), comment.UserID)
		if err == nil {
			followCount := uint(len(author.Follows))     // 统计关注数
			followerCount := uint(len(author.Followers)) // 统计粉丝数
			workCount := uint(len(author.Works))         // 统计作品数
			favoriteCount := uint(len(author.Favorites)) // 统计点赞数

			// 统计获赞数
			var favoritedCount uint = 0
			for _, video := range author.Works {
				favoritedCount += uint(len(video.Favorited))
			}

			// 是否关注
			isFollow := dao.CheckFollow(context.TODO(), Me_ID.(uint), author.ID)

			userInfo := response.User{
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

			// 评论信息中加入作者信息
			commentInfo.User = userInfo
		}
		// 将该评论加入列表
		resp.Comment_List = append(resp.Comment_List, commentInfo)
	}

	return resp, nil
}
