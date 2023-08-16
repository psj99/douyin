package service

import (
	"douyin/service/types/request"
	"douyin/service/types/response"

	"github.com/gin-gonic/gin"
)

func Comment(ctx *gin.Context, req *request.CommentReq) (resp *response.CommentResp, err error) {
	return &response.CommentResp{}, nil
}

func CommentList(ctx *gin.Context, req *request.CommentListReq) (resp *response.CommentListResp, err error) {
	return &response.CommentListResp{}, nil
}
