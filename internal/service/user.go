package service

import (
	"context"
	"douyin/internal/model"
	"douyin/internal/pkg/request"
	"douyin/internal/repository"
)

type UserService interface {
	Register(ctx context.Context, req *request.RegisterRequest) error
	Login(ctx context.Context, req *request.LoginRequest) (string, error)
	GetProfile(ctx context.Context, userId string) (*model.User, error)
	UpdateProfile(ctx context.Context, userId string, req *request.UpdateProfileRequest) error
}

type userService struct {
	*Service
	userRepo repository.UserRepository
}

func NewUserService(service *Service, userRepo repository.UserRepository) UserService {
	return &userService{
		Service:  service,
		userRepo: userRepo,
	}
}

func (s *userService) Register(ctx context.Context, req *request.RegisterRequest) error {
	//TODO implement me
	panic("implement me")
}

func (s *userService) Login(ctx context.Context, req *request.LoginRequest) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (s *userService) GetProfile(ctx context.Context, userId string) (*model.User, error) {
	//TODO implement me
	panic("implement me")
}

func (s *userService) UpdateProfile(ctx context.Context, userId string, req *request.UpdateProfileRequest) error {
	// TODO implement me
	panic("implement me")
}
