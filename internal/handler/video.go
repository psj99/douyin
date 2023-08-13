package handler

import (
	"douyin/internal/pkg/request"
	"douyin/internal/service"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"time"
)

type VideoHandler interface {
	GetFeed(ctx *gin.Context)
}

type videoHandler struct {
	*Handler
	videoService service.VideoService
}

func (v videoHandler) GetFeed(ctx *gin.Context) {
	var feedReq *request.FeedReq
	if err := ctx.ShouldBind(feedReq); err != nil {
		v.logger.Error("ShouldBind err", zap.Error(err))
	}
	if feedReq.LatestTime == 0 {
		feedReq.LatestTime = time.Now().Unix()
	}
	
	//TODO implement me
	panic("implement me")
}

func NewVideoHandler(handler *Handler, videoService service.VideoService) VideoHandler {
	return &videoHandler{
		Handler:      handler,
		videoService: videoService,
	}
}
