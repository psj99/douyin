package service

import (
	"context"
	"douyin/pkg/types/request"
	"douyin/pkg/types/response"
	"douyin/pkg/utils"
	"douyin/repository/dao"
	"douyin/repository/model"
	"errors"
	"gorm.io/gorm"
	"sync"
)

var userSrvInstance *userSrv
var userSrcOnce sync.Once

type userSrv struct {
}

func GetUserSrv() *userSrv {
	userSrcOnce.Do(func() {
		userSrvInstance = &userSrv{}
	})
	return userSrvInstance
}

func (this *userSrv) UserRegister(ctx context.Context, req *request.UserRegisterReq) (resp any, err error) {
	userDao := dao.NewUserDao(ctx)
	_, err = userDao.FindUserByUserName(req.Username)
	if err == nil {
		err = errors.New("用户已存在")
		return
	}
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		utils.ZapLogger.Infof("FindUserByUserName err: %v", err)
		return
	}
	//只有err = ErrRecordNotFound时才可以注册
	user := &model.User{
		UserName: req.Username,
	}
	err = user.SetPassword(req.Password)
	if err != nil {
		utils.ZapLogger.Infof("SetPassword err %v", err)
		return
	}
	err = userDao.CreateUser(user)
	if err != nil {
		utils.ZapLogger.Infof("CreateUser err %v", err)
		return
	}

	userLoginReq := &request.UserLoginReq{
		Username: req.Username,
		Password: req.Password,
	}
	resp, err = this.UserLogin(ctx, userLoginReq)
	return resp, nil
}

func (this *userSrv) UserLogin(ctx context.Context, req *request.UserLoginReq) (resp any, err error) {
	userDao := dao.NewUserDao(ctx)
	user, err := userDao.FindUserByUserName(req.Username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = errors.New("用户不存在")
			return nil, err
		}
		utils.ZapLogger.Infof("FindUserByUserName err: %v", err)
		return nil, err
	}
	if !user.CheckPassword(req.Password) {
		err = errors.New("账号或密码错误")
		return nil, err
	}
	token, err := utils.GenerateToken(user.ID, user.UserName)
	if err != nil {
		utils.ZapLogger.Infof("GenerateToken err: %v", err)
		err = errors.New("生成token错误")
		return nil, err
	}
	resp = &response.UserLoginResp{
		Code:   0,
		Msg:    "成功",
		UserId: user.ID,
		Token:  token,
	}
	return resp, nil

}
