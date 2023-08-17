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
	req_id, ok := ctx.Get("user_id")
	if !ok {
		utils.ZapLogger.Errorf("ctx.Get (user_id) err: inaccessible")
		return nil, errors.New("无法获取请求用户ID")
	}

	// 读取目标视频ID
	video_id, err := strconv.ParseUint(req.Video_ID, 10, 64)
	if err != nil {
		utils.ZapLogger.Errorf("ParseUint err: %v", err)
		return nil, err
	}

	// 存储评论信息
	action_type, err := strconv.ParseUint(req.Action_Type, 10, 64)
	if err != nil {
		utils.ZapLogger.Errorf("ParseUint err: %v", err)
		return nil, err
	}
	if action_type == 1 {
		// 创建评论
		// 存储评论信息 //TODO
		comment, err := dao.CreateComment(context.TODO(), req_id.(uint), uint(video_id), req.Comment_Text)
		if err != nil {
			utils.ZapLogger.Errorf("CreateComment err: %v", err)
			return nil, err
		}

		// 初始化响应
		commentInfo := response.Comment{ID: comment.ID, Content: comment.Content}

		// 评论发布时间
		commentInfo.Create_Date = fmt.Sprintf("%02d-%02d", comment.CreatedAt.Month(), comment.CreatedAt.Day()) // mm-dd

		// 临时方案 亟待优化 //TODO
		// 评论作者信息
		authorInfo, err := readUserInfo(ctx, req_id.(uint))
		if err != nil {
			// 响应为评论成功 但作者将为空
			utils.ZapLogger.Errorf("readUserInfo err: %v", err)
		} else {
			commentInfo.User = *authorInfo
		}

		return &response.CommentResp{Comment: commentInfo}, nil
	} else if action_type == 2 {
		// 删除评论
		// 读取目标评论ID
		comment_id, err := strconv.ParseUint(req.Comment_ID, 10, 64)
		if err != nil {
			utils.ZapLogger.Errorf("ParseUint err: %v", err)
			return nil, err
		}

		// 删除评论信息
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
	// 读取目标视频ID
	video_id, err := strconv.ParseUint(req.Video_ID, 10, 64)
	if err != nil {
		utils.ZapLogger.Errorf("ParseUint err: %v", err)
		return nil, err
	}

	// 读取目标视频评论列表
	comments, err := dao.FindCommentsByCreatedAt(context.TODO(), uint(video_id), false)
	if err != nil {
		utils.ZapLogger.Errorf("FindCommentsByCreatedAt err: %v", err)
		return nil, err
	}

	resp = &response.CommentListResp{} // 初始化响应
	for _, comment := range comments {
		// 读取评论信息
		commentInfo := response.Comment{ID: comment.ID, Content: comment.Content}

		// 评论发布时间
		commentInfo.Create_Date = fmt.Sprintf("%02d-%02d", comment.CreatedAt.Month(), comment.CreatedAt.Day()) // mm-dd

		// 读取作者信息
		authorInfo, err := readUserInfo(ctx, comment.UserID)
		if err != nil {
			utils.ZapLogger.Errorf("readUserInfo err: %v", err)
			continue // 跳过本条评论
		} else {
			commentInfo.User = *authorInfo
		}

		// 将该评论加入列表
		resp.Comment_List = append(resp.Comment_List, commentInfo)
	}

	return resp, nil
}
