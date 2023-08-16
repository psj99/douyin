package service

import (
	"douyin/repository/dao"
	"douyin/service/types/request"
	"douyin/service/types/response"
	"douyin/utils"

	"context"
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
)

func Message(ctx *gin.Context, req *request.MessageReq) (resp *response.MessageResp, err error) {
	// 获取请求用户ID
	Me_ID, ok := ctx.Get("user_id")
	if !ok {
		utils.ZapLogger.Errorf("ctx.Get (user_id) err: inaccessible")
		return nil, errors.New("无法获取user_id")
	}

	// 获取目标用户ID
	to_user_id, err := strconv.ParseUint(req.To_User_ID, 10, 64) // string转十进制uint64
	if err != nil {
		utils.ZapLogger.Errorf("ParseUint err: %v", err)
		return nil, err
	}

	// 发送消息
	action_type, err := strconv.ParseUint(req.Action_Type, 10, 64)
	if err != nil {
		utils.ZapLogger.Errorf("ParseUint err: %v", err)
		return nil, err
	}
	if action_type == 1 {
		_, err := dao.CreateMessage(context.TODO(), Me_ID.(uint), uint(to_user_id), req.Content)
		if err != nil {
			utils.ZapLogger.Errorf("CreateMessage err: %v", err)
			return nil, err
		}
	} else {
		utils.ZapLogger.Errorf("Invalid action_type err: %v", action_type)
		return nil, errors.New("操作类型有误")
	}

	return &response.MessageResp{}, nil
}

func MessageList(ctx *gin.Context, req *request.MessageListReq) (resp *response.MessageListResp, err error) {
	// 获取请求用户ID
	Me_ID, ok := ctx.Get("user_id")
	if !ok {
		utils.ZapLogger.Errorf("ctx.Get (user_id) err: inaccessible")
		return nil, errors.New("无法获取user_id")
	}

	// 获取目标用户ID
	to_user_id, err := strconv.ParseUint(req.To_User_ID, 10, 64) // string转十进制uint64
	if err != nil {
		utils.ZapLogger.Errorf("ParseUint err: %v", err)
		return nil, err
	}

	// 读取消息列表
	outMessages, err1 := dao.FindMessagesBy_From_To_ID(context.TODO(), Me_ID.(uint), uint(to_user_id), int64(req.Pre_Msg_Time), true, -1) // 查找从某刻起新消息 不限制数量
	if err1 != nil {
		utils.ZapLogger.Errorf("FindMessagesBy_From_To_ID err: %v", err)
		return nil, err1
	}
	inMessages, err2 := dao.FindMessagesBy_From_To_ID(context.TODO(), uint(to_user_id), Me_ID.(uint), int64(req.Pre_Msg_Time), true, -1) // 查找从某刻起新消息 不限制数量
	if err2 != nil {
		utils.ZapLogger.Errorf("FindMessagesBy_From_To_ID err: %v", err)
		return nil, err2
	}
	messages := append(outMessages, inMessages...)

	resp = &response.MessageListResp{}
	for _, message := range messages {
		messageInfo := response.Message{
			ID:           message.ID,
			To_User_ID:   message.ToUserID,
			From_User_ID: message.FromUserID,
			Content:      message.Content,
			Create_Time:  uint(message.CreatedAt.Unix() * 1000), // 消息发送时间 API文档有误 实为毫秒数时间戳
			// Create_Time:  message.CreatedAt.Format("2006-01-02 15:04:05"),
		}

		resp.Message_List = append(resp.Message_List, messageInfo)
	}

	return resp, nil
}
