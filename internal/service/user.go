package service

import (
	"context"
	"douyin/internal/model"
	"douyin/internal/pkg/request"
	"douyin/internal/pkg/resp"
	"douyin/internal/repository"
	"errors"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"strconv"
)

type UserService interface {
	//Register(ctx context.Context, req *request.UserRegisterReq) (*resp.UserRegisterResp, error)
	Register(ctx context.Context, req *request.UserRegisterReq) (userId uint, token string, err error)
	//Login(ctx context.Context, req *request.UserLoginReq) (*resp.UserLoginResp, error)
	Login(ctx context.Context, req *request.UserLoginReq) (userId uint, token string, err error)
	GetUserInfo(ctx context.Context, userId uint) (*resp.UserInfo, error)
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

func (s *userService) Register(ctx context.Context, req *request.UserRegisterReq) (userId uint, token string, err error) {
	// 查找用户是否已存在
	_, err = s.userRepo.FindUserByUserName(ctx, req.Username)
	if err == nil {
		err = errors.New("用户名已存在")
		return 0, "", err
	} else if !errors.Is(err, gorm.ErrRecordNotFound) { // 若出现查找功能性错误
		s.logger.Error("FindUserByUserName err:", zap.Error(err))
		return 0, "", err
	}

	// 只有err == ErrRecordNotFound时才可以注册
	// 准备要存储的内容
	user := &model.User{
		UserName: req.Username,
	}
	err = user.SetPassword(req.Password) // 向要存储的内容中添加单向加密后的密码
	if err != nil {
		s.logger.Error("SetPassword err:", zap.Error(err))
		return 0, "", err
	}

	// 存储用户信息
	err = s.userRepo.CreateUser(ctx, user)
	if err != nil {
		s.logger.Error("CreateUser err:", zap.Error(err))
		return 0, "", err
	}

	// 注册后自动登录
	userLoginReq := &request.UserLoginReq{
		Username: req.Username,
		Password: req.Password,
	}
	userId, token, err = s.Login(ctx, userLoginReq)
	if err != nil {
		s.logger.Error("UserLogin err:", zap.Error(err))
		return 0, "", err
	}

	return userId, token, nil
	//return &response.UserRegisterResp{StatusCode: 0, StatusMsg: "注册成功", UserId: loginResp.UserId, Token: loginResp.Token}, err
}

func (s *userService) Login(ctx context.Context, req *request.UserLoginReq) (userId uint, token string, err error) {
	// 查找用户是否已存在
	user, err := s.userRepo.FindUserByUserName(ctx, req.Username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = errors.New("用户不存在")
			return 0, "", err
		} else { // 若出现查找功能性错误
			s.logger.Error("FindUserByUserName err:", zap.Error(err))
			return 0, "", err
		}
	}

	// 校验密码
	if !user.CheckPassword(req.Password) {
		err = errors.New("用户名或密码错误")
		return 0, "", err
	}

	// 校验成功时生成用户鉴权token
	token, err = s.jwt.GenToken(strconv.Itoa(int(user.ID)))
	if err != nil {
		s.logger.Error("GenerateToken err:", zap.Error(err))
		return 0, "", err
	}

	return user.ID, token, err

}

func (s *userService) GetUserInfo(ctx context.Context, userId uint) (*resp.UserInfo, error) {
	user, err := s.userRepo.FindUserByUserID(ctx, userId)
	if err != nil {
		s.logger.Error("FindUserByUserID err", zap.Error(err))
		return nil, err
	}
	userInfo := &resp.UserInfo{
		UserId:          int64(user.ID),
		Name:            user.UserName,
		FollowCount:     0,
		FollowerCount:   0,
		IsFollow:        false,
		Avatar:          user.Avatar,
		BackGroundImage: user.BackgroundImage,
		Signature:       user.Signature,
		TotalFavorited:  "",
		WorkCount:       "",
		FavoriteCount:   "",
	}
	return userInfo, nil
}
