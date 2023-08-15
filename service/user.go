package service

import (
	"douyin/repository/dao"
	"douyin/repository/model"
	"douyin/service/types/request"
	"douyin/service/types/response"
	"douyin/utils"

	"context"
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// 自定义错误类型
var ErrorUserExists = errors.New("用户已存在")
var ErrorUserNotExists = errors.New("用户不存在")
var ErrorWrongPassword = errors.New("账号或密码错误")

// 用户注册
func UserRegister(ctx *gin.Context, req *request.UserRegisterReq) (resp *response.UserRegisterResp, err error) {
	// 查找用户是否已存在
	_, err = dao.FindUserByUsername(context.TODO(), req.Username)
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
		Username: req.Username,
	}
	err = user.SetPassword(req.Password) // 向要存储的内容中添加单向加密后的密码
	if err != nil {
		utils.ZapLogger.Errorf("SetPassword err: %v", err)
		return nil, err
	}

	// 存储用户信息
	user, err = dao.CreateUser(context.TODO(), user)
	if err != nil {
		utils.ZapLogger.Errorf("CreateUser err: %v", err)
		return nil, err
	}

	// 注册后生成用户鉴权token(自动登录)
	token, err := utils.GenerateToken(user.ID, user.Username)
	if err != nil {
		utils.ZapLogger.Errorf("GenerateToken err: %v", err)
		return nil, err
	}

	return &response.UserRegisterResp{User_Id: user.ID, Token: token}, err
}

// 用户登录
func UserLogin(ctx *gin.Context, req *request.UserLoginReq) (resp *response.UserLoginResp, err error) {
	// 查找用户是否已存在
	user, err := dao.FindUserByUsername(context.TODO(), req.Username)
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
	token, err := utils.GenerateToken(user.ID, user.Username)
	if err != nil {
		utils.ZapLogger.Errorf("GenerateToken err: %v", err)
		return nil, err
	}

	return &response.UserLoginResp{User_Id: user.ID, Token: token}, err
}

// 用户信息
func UserInfo(ctx *gin.Context, req *request.UserInfoReq) (resp *response.UserInfoResp, err error) {
	// 读取用户信息
	userID, err := strconv.ParseUint(req.User_ID, 10, 64) // string转十进制uint64
	if err != nil {
		utils.ZapLogger.Errorf("ParseUint err: %v", err)
		return nil, err
	}
	user, err := dao.FindUserByID(context.TODO(), uint(userID))
	if err != nil {
		utils.ZapLogger.Errorf("FindUserByID err: %v", err)
		return nil, err
	}

	// 临时方案 亟待优化 (需要大幅更改数据库模型) //TODO
	followCount := uint(len(user.Follows))     // 统计关注数
	followerCount := uint(len(user.Followers)) // 统计粉丝数
	workCount := uint(len(user.Works))         // 统计作品数
	favoriteCount := uint(len(user.Favorites)) // 统计点赞数

	// 统计获赞数
	var favoritedCount uint = 0
	for _, video := range user.Works {
		favoritedCount += uint(len(video.Favorited))
	}

	// 是否关注
	Me_ID, ok := ctx.Get("user_id") // 获取自己
	if !ok {
		utils.ZapLogger.Errorf("ctx.Get (user_id) err: %v", err)
		return nil, errors.New("无法获取user_id")
	}
	isFollow, err := dao.CheckFollower(context.TODO(), uint(userID), Me_ID.(uint))
	if err != nil {
		isFollow = false
		utils.ZapLogger.Errorf("CheckFollow err: %v", err)
	}

	return &response.UserInfoResp{User: response.User{
		ID:               user.ID,
		Name:             user.Username,
		Follow_Count:     followCount,
		Follower_Count:   followerCount,
		Is_Follow:        isFollow,
		Avatar:           user.Avatar,
		Background_Image: user.BackgroundImage,
		Signature:        user.Signature,
		Total_Favorited:  strconv.FormatUint(uint64(favoritedCount), 10),
		Work_Count:       workCount,
		Favorite_Count:   favoriteCount,
	}}, err
}
