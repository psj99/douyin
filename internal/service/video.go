package service

import (
	"douyin/internal/model"
	"douyin/internal/repository"
)

type VideoService interface {
	GetVideoById(id int64) (*model.Video, error)
}

type videoService struct {
	*Service
	videoRepository repository.VideoRepository
}

func NewVideoService(service *Service, videoRepository repository.VideoRepository) VideoService {
	return &videoService{
		Service:        service,
		videoRepository: videoRepository,
	}
}

func (s *videoService) GetVideoById(id int64) (*model.Video, error) {
	return s.videoRepository.FirstById(id)
}
