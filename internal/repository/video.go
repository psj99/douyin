package repository

import (
	"douyin/internal/model"
)

type VideoRepository interface {
	FirstById(id int64) (*model.Video, error)
}
type videoRepository struct {
	*Repository
}

func NewVideoRepository(repository *Repository) VideoRepository {
	return &videoRepository{
		Repository: repository,
	}
}

func (r *videoRepository) FirstById(id int64) (*model.Video, error) {
	var video model.Video
	// TODO: query db
	return &video, nil
}
