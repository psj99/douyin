package dao

import (
	"douyin/repository/model"

	"context"

	"gorm.io/gorm/clause"
)

// 根据用户id找到用户
func FindUserByID(ctx context.Context, id uint) (user *model.User, err error) {
	DB := GetDB(ctx)
	user = &model.User{}
	err = DB.Model(&model.User{}).Where("id=?", id).Preload(clause.Associations).First(user).Error
	return user, err
}

// 根据用户名找到用户
func FindUserByUsername(ctx context.Context, username string) (user *model.User, err error) {
	DB := GetDB(ctx)
	user = &model.User{}
	err = DB.Model(&model.User{}).Where("username=?", username).Preload(clause.Associations).First(user).Error
	return user, err
}

// 创建User
func CreateUser(ctx context.Context, userInfo *model.User) (user *model.User, err error) {
	DB := GetDB(ctx)
	err = DB.Model(&model.User{}).Create(userInfo).Error
	return userInfo, err
}

// 检查是否关注
func CheckFollower(ctx context.Context, id uint, follower_id uint) (is_follower bool, err error) {
	DB := GetDB(ctx)
	var count int64
	err = DB.Model(&model.User{}).Preload("Followers").Where("id=? AND follower_id=?", id, follower_id).Count(&count).Error
	return count > 0, err
}

// 检查是否点赞
func CheckFavorite(ctx context.Context, id uint, video_id uint) (is_favorite bool, err error) {
	DB := GetDB(ctx)
	var count int64
	err = DB.Model(&model.User{}).Preload("Favorites").Where("id=? AND favorites_id=?", id, video_id).Count(&count).Error
	return count > 0, err
}
