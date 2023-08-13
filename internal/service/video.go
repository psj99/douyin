package service

import (
	"context"
	"douyin/internal/model"
	"douyin/internal/repository"
)

type VideoService interface {
	GetVideoById(ctx context.Context, id int64) (*model.Video, error)
	GetVideoList(ctx context.Context, latestTime int64, token string) ([]*model.Video, error)
}

type videoService struct {
	*Service
	videoRepository repository.VideoRepository
}

func (s *videoService) GetVideoById(ctx context.Context, id int64) (*model.Video, error) {
	//TODO implement me
	panic("implement me")
}

func (s *videoService) GetVideoList(ctx context.Context, latestTime int64, token string) ([]*model.Video, error) {
	//TODO implement me
	panic("implement me")
}

func NewVideoService(service *Service, videoRepository repository.VideoRepository) VideoService {
	return &videoService{
		Service:         service,
		videoRepository: videoRepository,
	}
}
