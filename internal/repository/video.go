package repository

import (
	"context"
	"douyin/internal/model"
)

type VideoRepository interface {
	FirstById(ctx context.Context, id int64) (*model.Video, error)
}
type videoRepository struct {
	*Repository
}

func (r *videoRepository) FirstById(ctx context.Context, id int64) (*model.Video, error) {
	//TODO implement me
	panic("implement me")
}

func NewVideoRepository(repository *Repository) VideoRepository {
	return &videoRepository{
		Repository: repository,
	}
}
