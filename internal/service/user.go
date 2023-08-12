package service

import (
	"context"
	"douyin/internal/pkg/request"
	"douyin/internal/pkg/response"
	"douyin/internal/repository"
)

type UserService interface {
	Register(ctx context.Context, req *request.UserRegisterReq) (*response.UserRegisterResp, error)
	Login(ctx context.Context, req *request.UserLoginReq) (*response.UserLoginResp, error)
	GetUserInfo(ctx context.Context, userId uint) (*response.UserInfo, error)
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

func (s *userService) Register(ctx context.Context, req *request.UserRegisterReq) (*response.UserRegisterResp, error) {
	//TODO implement me
	panic("implement me")
}

func (s *userService) Login(ctx context.Context, req *request.UserLoginReq) (*response.UserLoginResp, error) {
	//TODO implement me
	panic("implement me")
}

func (s *userService) GetUserInfo(ctx context.Context, userId uint) (*response.UserInfo, error) {
	//TODO implement me
	panic("implement me")
}
