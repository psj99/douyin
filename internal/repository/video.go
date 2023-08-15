package repository

import (
	"context"
	"douyin/internal/model"
	"fmt"
	"time"
)

type VideoRepository interface {
	Create(ctx context.Context, video *model.Video) error
	GetVideoList(ctx context.Context, latest int64) (videoList []*model.Video, err error)
	GetVideoListById(ctx context.Context, userId uint) (videoList []*model.Video, err error)
}
type videoRepository struct {
	*Repository
}

func (vRepo *videoRepository) GetVideoListById(ctx context.Context, userId uint) (videoList []*model.Video, err error) {
	err = vRepo.db.Model(&model.Video{}).Where("user_id = ?", userId).Order("create_at desc").Find(videoList).Error
	return
}

func (vRepo *videoRepository) GetVideoList(ctx context.Context, latest int64) (videoList []*model.Video, err error) {
	latestTime := fmt.Sprintf(time.Unix(latest, 0).Format("2006-01-02 15:04:05"))
	err = vRepo.db.Model(&model.Video{}).Where("create_at < ?", latestTime).Order("create_at desc").Find(videoList).Error
	if err != nil {
		// 交给service层处理，不需要记录日志
		//vRepo.logger.Error("GetVideoList err:", zap.Error(err))
		return nil, err
	}
	return videoList, nil

}

func (vRepo *videoRepository) Create(ctx context.Context, video *model.Video) error {
	err := vRepo.db.Model(&model.User{}).Create(video).Error
	return err
}

func NewVideoRepository(repository *Repository) VideoRepository {
	return &videoRepository{
		Repository: repository,
	}
}
