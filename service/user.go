package service

import (
	"douyin/repository/dao"
	"douyin/repository/model"
	"douyin/service/types/request"
	"douyin/service/types/response"
	"douyin/utils"

	"context"
	"errors"

	"gorm.io/gorm"
)

// 自定义错误类型
var ErrorUserExists = errors.New("用户已存在")
var ErrorUserNotExists = errors.New("用户不存在")
var ErrorWrongPassword = errors.New("账号或密码错误")

// 用户注册
func UserRegister(ctx context.Context, req *request.UserRegisterReq) (resp *response.UserRegisterResp, err error) {
	// 查找用户是否已存在
	_, err = dao.FindUserByUserName(ctx, req.Username)
	if err == nil {
		err = ErrorUserExists
		return nil, err
	} else if !errors.Is(err, gorm.ErrRecordNotFound) { // 若出现查找功能性错误
		utils.ZapLogger.Errorf("FindUserByUserName err: %v", err)
		return nil, err
	}

	// 只有err == ErrRecordNotFound时才可以注册
	// 准备要存储的内容
	user := &model.User{
		UserName: req.Username,
	}
	err = user.SetPassword(req.Password) // 向要存储的内容中添加单向加密后的密码
	if err != nil {
		utils.ZapLogger.Errorf("SetPassword err %v", err)
		return nil, err
	}

	// 存储用户信息
	err = dao.CreateUser(ctx, user)
	if err != nil {
		utils.ZapLogger.Errorf("CreateUser err %v", err)
		return nil, err
	}

	// 注册后自动登录
	userLoginReq := &request.UserLoginReq{
		Username: req.Username,
		Password: req.Password,
	}
	loginResp, err := UserLogin(ctx, userLoginReq)
	if err != nil {
		utils.ZapLogger.Errorf("UserLogin err %v", err)
		return nil, err
	}

	return &response.UserRegisterResp{Status_Code: 0, Status_Msg: "注册成功", User_Id: loginResp.User_Id, Token: loginResp.Token}, err
}

// 用户登录
func UserLogin(ctx context.Context, req *request.UserLoginReq) (resp *response.UserLoginResp, err error) {
	// 查找用户是否已存在
	user, err := dao.FindUserByUserName(ctx, req.Username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = ErrorUserNotExists
			return nil, err
		} else { // 若出现查找功能性错误
			utils.ZapLogger.Errorf("FindUserByUserName err: %v", err)
			return nil, err
		}
	}

	// 校验密码
	if !user.CheckPassword(req.Password) {
		err = ErrorWrongPassword
		return nil, err
	}

	// 校验成功时生成用户鉴权token
	token, err := utils.GenerateToken(user.ID, user.UserName)
	if err != nil {
		utils.ZapLogger.Errorf("GenerateToken err: %v", err)
		return nil, err
	}

	return &response.UserLoginResp{Status_Code: 0, Status_Msg: "登录成功", User_Id: user.ID, Token: token}, err
}
