package repository

import (
	"context"
	"douyin/internal/model"
)

type UserRepository interface {
	FindUserByUserID(ctx context.Context, id uint) (*model.User, error)
	FindUserByUserName(ctx context.Context, username string) (*model.User, error)
	CreateUser(ctx context.Context, user *model.User) error
}

type userRepository struct {
	*Repository
}

func (u userRepository) FindUserByUserID(ctx context.Context, id uint) (*model.User, error) {
	user := &model.User{}
	err := u.db.Model(&model.User{}).Where("id=?", id).First(user).Error
	return user, err
}

func (u userRepository) FindUserByUserName(ctx context.Context, username string) (*model.User, error) {
	user := &model.User{}
	err := u.db.Model(&model.User{}).Where("user_name=?", username).First(user).Error
	return user, err
}

func (u userRepository) CreateUser(ctx context.Context, user *model.User) error {
	err := u.db.Model(&model.User{}).Create(user).Error
	return err
}

func NewUserRepository(r *Repository) UserRepository {
	return &userRepository{
		Repository: r,
	}
}
