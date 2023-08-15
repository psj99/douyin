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
func CheckFollow(ctx context.Context, id uint, follow_id uint) (is_follower bool, err error) {
	if id == follow_id {
		return false, nil // 默认自己不关注自己
	}

	DB := GetDB(ctx)
	var count int64
	err = DB.Table("follow").Where("user_id=? AND follow_id=?", id, follow_id).Count(&count).Error
	if err != nil {
		return false, err
	} else {
		return count > 0, nil
	}
}

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
