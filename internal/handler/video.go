package handler

import (
	"douyin/internal/service"
	"douyin/pkg/helper/qiniu"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

type VideoHandler interface {
	GetFeed(ctx *gin.Context)
	Upload(ctx *gin.Context)
}

type videoHandler struct {
	*Handler
	videoService service.VideoService
}

func (v videoHandler) Upload(ctx *gin.Context) {
	// 解析上传的文件
	file, err := ctx.FormFile("file")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	uploader := qiniu.NewQiniuUploader(
		"s3RYVO1nDvkpx8GFOgzySq_nRp7hefFNkF2QFRvj",
		"oVDq14H6LrMwkeBwfjS-1adlDDfPbyTdv5J80K7a",
		"tk-repo",
		"ryv7jqdrm.hn-bkt.clouddn.com")

	fileURL, coverURL, err := uploader.UploadFile(ctx, file)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"err": err.Error()})
		return
	}
	v.logger.Info("上传成功", zap.String("fileURL", fileURL), zap.String("coverURL", fileURL))
	ctx.JSON(http.StatusOK, gin.H{"fileURL": fileURL, "coverURL": coverURL})
}

func (v videoHandler) GetFeed(ctx *gin.Context) {
	//TODO implement me
	panic("implement me")
}

func NewVideoHandler(handler *Handler, videoService service.VideoService) VideoHandler {
	return &videoHandler{
		Handler:      handler,
		videoService: videoService,
	}
}
