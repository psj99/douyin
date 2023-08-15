package handler

import (
	"douyin/internal/pkg/resp"
	"douyin/internal/service"
	"douyin/pkg/helper/convert"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type VideoHandler interface {
	Feed(ctx *gin.Context)
	PublishAction(ctx *gin.Context)
	PublishList(ctx *gin.Context)
}

type videoHandler struct {
	*Handler
	videoService service.VideoService
	userService  service.UserService
}

func (videoHandler videoHandler) PublishList(ctx *gin.Context) {
	userId := convert.StringToUint(GetUserIdFromCtx(ctx))
	videoes, err := videoHandler.videoService.GetPublish(ctx, userId)
	if err != nil {
		return
	}

	user, _ := videoHandler.userService.GetUserInfo(ctx, userId)

	var videoList []*resp.Video
	for _, video := range videoes {
		videoInfo := &resp.Video{
			Id:            int64(video.ID),
			User:          user,
			PlayUrl:       video.PlayUrl,
			CoverUrl:      video.CoverUrl,
			FavoriteCount: 0,
			CommentCount:  0,
			IsFavorite:    false,
			Title:         video.Title,
		}
		videoList = append(videoList, videoInfo)
	}
	ctx.JSON(http.StatusOK, &resp.PublishListResp{
		Response:  resp.ResponseOK(),
		VideoList: videoList,
	})

}

func (videoHandler videoHandler) PublishAction(ctx *gin.Context) {

	title := ctx.Query("title")
	// 解析上传的文件
	file, err := ctx.FormFile("file")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userId := convert.StringToUint(GetUserIdFromCtx(ctx))
	err = videoHandler.videoService.PublishVideo(ctx, file, uint(userId), title)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, resp.PublishActionResp{Response: resp.ResponseErr(err.Error())})
		return
	}

	ctx.JSON(http.StatusOK, resp.PublishActionResp{
		Response: resp.ResponseOK(),
	})

	//uploader := qiniu.NewQiniuUploader(
	//	"s3RYVO1nDvkpx8GFOgzySq_nRp7hefFNkF2QFRvj",
	//	"oVDq14H6LrMwkeBwfjS-1adlDDfPbyTdv5J80K7a",
	//	"tk-repo",
	//	"ryv7jqdrm.hn-bkt.clouddn.com")
	//
	//fileURL, coverURL, err := uploader.UploadFile(ctx, file)
	//if err != nil {
	//	ctx.JSON(http.StatusOK, gin.H{"err": err.Error()})
	//	return
	//}
	//videoHandler.logger.Info("上传成功", zap.String("fileURL", fileURL), zap.String("coverURL", fileURL))
	//ctx.JSON(http.StatusOK, gin.H{"fileURL": fileURL, "coverURL": coverURL})
}

func (videoHandler videoHandler) Feed(ctx *gin.Context) {

	latest := ctx.Query("latest_time")
	var latestUnix int64 = 0
	if latest == "" {
		latestUnix = time.Now().Unix()
	} else {
		latestUnix = convert.StringToInt64(latest)
	}
	videoes, nextTime, err := videoHandler.videoService.GetFeed(ctx, latestUnix)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, resp.FeedResp{
			Response:  resp.ResponseErr(err.Error()),
			VideoList: nil,
			NextTime:  nextTime,
		})
		return
	}

	var videoList []*resp.Video
	for _, video := range videoes {
		user, _ := videoHandler.userService.GetUserInfo(ctx, video.UserID)
		videoInfo := &resp.Video{
			Id:            int64(video.ID),
			User:          user,
			PlayUrl:       video.PlayUrl,
			CoverUrl:      video.CoverUrl,
			FavoriteCount: 0,
			CommentCount:  0,
			IsFavorite:    false,
			Title:         video.Title,
		}
		videoList = append(videoList, videoInfo)
	}

	ctx.JSON(http.StatusOK, resp.FeedResp{
		Response:  resp.ResponseOK(),
		VideoList: videoList,
		NextTime:  nextTime,
	})
}

func NewVideoHandler(handler *Handler, videoService service.VideoService) VideoHandler {
	return &videoHandler{
		Handler:      handler,
		videoService: videoService,
	}
}
