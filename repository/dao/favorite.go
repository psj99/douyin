package dao

import (
	"douyin/repository/model"

	"context"
)

// 点赞
func CreateFavorite(ctx context.Context, user_id uint, video_id uint) (err error) {
	DB := GetDB(ctx)
	user := &model.User{}
	err = DB.Model(&model.User{}).Where("id=?", user_id).First(user).Error
	if err != nil {
		return err
	}
	video := &model.Video{}
	err = DB.Model(&model.Video{}).Where("id=?", video_id).First(video).Error
	if err != nil {
		return err
	}
	return DB.Model(user).Association("Favorites").Append(video)
}

// 取消点赞
func DeleteFavorite(ctx context.Context, user_id uint, video_id uint) (err error) {
	DB := GetDB(ctx)
	user := &model.User{}
	err = DB.Model(&model.User{}).Where("id=?", user_id).First(user).Error
	if err != nil {
		return err
	}
	video := &model.Video{}
	err = DB.Model(&model.Video{}).Where("id=?", video_id).First(video).Error
	if err != nil {
		return err
	}
	return DB.Model(user).Association("Favorites").Delete(video)
}

// 检查是否点赞
func CheckFavorite(ctx context.Context, user_id uint, video_id uint) (is_favorite bool, err error) {
	DB := GetDB(ctx)
	var count int64
	err = DB.Table("favorite").Where("user_id=? AND video_id=?", user_id, video_id).Count(&count).Error
	if err != nil {
		return false, err
	} else {
		return count > 0, nil
	}
}
