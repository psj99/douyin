package service

import (
	"context"
	"douyin/internal/model"
	"douyin/internal/repository"
	"go.uber.org/zap"
	"mime/multipart"
)

type VideoService interface {
	PublishVideo(ctx context.Context, file *multipart.FileHeader, userId uint, title string) error
	GetFeed(ctx context.Context, latestUnix int64) (videoList []*model.Video, nextTime int64, err error)
	GetPublish(ctx context.Context, userId uint) (videoList []*model.Video, err error)
}

type videoService struct {
	*Service
	videoRepository repository.VideoRepository
}

func (vService *videoService) GetPublish(ctx context.Context, userId uint) (videoList []*model.Video, err error) {
	videoList, err = vService.videoRepository.GetVideoListById(ctx, userId)
	if err != nil {
		vService.logger.Error("GetVideoList发生了错误：", zap.Error(err))
		return nil, err
	}
	return
}

func (vService *videoService) GetFeed(ctx context.Context, latestUnix int64) (videoList []*model.Video, nextTime int64, err error) {

	videoList, err = vService.videoRepository.GetVideoList(ctx, latestUnix)
	if err != nil {
		vService.logger.Error("GetVideoList发生了错误：", zap.Error(err))
		return nil, latestUnix, err
	}

	nextTime = videoList[0].CreatedAt.Unix()
	return
}

func (vService *videoService) PublishVideo(ctx context.Context, file *multipart.FileHeader, userId uint, title string) error {
	//TODO implement me
	videoURL, coverURL, err := vService.uploader.UploadFile(ctx, file)
	if err != nil {
		vService.logger.Error("七牛云upload发生错误: ", zap.String("err", err.Error()))
		return err
	}
	newVideo := &model.Video{
		PlayUrl:    videoURL,
		CoverUrl:   coverURL,
		Title:      title,
		UserID:     userId,
		Comments:   nil,
		LikedUsers: nil,
	}
	err = vService.videoRepository.Create(ctx, newVideo)
	if err != nil {
		vService.logger.Error("向数据库添加Video发生错误: ", zap.String("err", err.Error()))
		return err
	}
	return nil
}

func NewVideoService(service *Service, videoRepository repository.VideoRepository) VideoService {
	return &videoService{
		Service:         service,
		videoRepository: videoRepository,
	}
}
