package dao

import (
	"douyin/repository/model"

	"context"
)

// 根据用户id找到用户
func FindUserByUserId(ctx context.Context, id uint) (user *model.User, err error) {
	DB := GetDB(ctx)
	user = &model.User{}
	err = DB.Model(&model.User{}).Where("id=?", id).First(user).Error
	return user, err
}

// 根据用户名找到用户
func FindUserByUserName(ctx context.Context, userName string) (user *model.User, err error) {
	DB := GetDB(ctx)
	user = &model.User{}
	err = DB.Model(&model.User{}).Where("user_name=?", userName).First(user).Error
	return user, err
}

// 创建User
func CreateUser(ctx context.Context, userInfo *model.User) (user *model.User, err error) {
	DB := GetDB(ctx)
	err = DB.Model(&model.User{}).Create(userInfo).Error
	return userInfo, err
}
