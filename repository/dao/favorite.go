package dao

import (
	"douyin/repository/model"

	"context"
)

// 检查是否点赞
func CheckFavorite(ctx context.Context, id uint, video_id uint) (is_favorite bool, err error) {
	DB := GetDB(ctx)
	var count int64
	err = DB.Table("favorite").Where("user_id=? AND video_id=?", id, video_id).Count(&count).Error
	if err != nil {
		return false, err
	} else {
		return count > 0, nil
	}
}

// 点赞
func CreateFavorite(ctx context.Context, id uint, video_id uint) (err error) {
	DB := GetDB(ctx)
	video := &model.Video{}
	err = DB.Model(&model.Video{}).Where("id=?", video_id).First(video).Error
	if err != nil {
		return err
	}
	return DB.Model(&model.User{}).Where("id=?", id).Association("Favorites").Append(video)
}

// 取消点赞
func DeleteFavorite(ctx context.Context, id uint, video_id uint) (err error) {
	DB := GetDB(ctx)
	video := &model.Video{}
	err = DB.Model(&model.Video{}).Where("id=?", video_id).First(video).Error
	if err != nil {
		return err
	}
	return DB.Model(&model.User{}).Where("id=?", id).Association("Favorites").Delete(video)
}
