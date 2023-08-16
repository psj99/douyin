package dao

import (
	"douyin/repository/model"

	"context"

	"gorm.io/gorm/clause"
)

// 发送消息
func CreateMessage(ctx context.Context, from_user_id uint, to_user_id uint, content string) (message *model.Message, err error) {
	DB := GetDB(ctx)
	message = &model.Message{Content: content, FromUserID: from_user_id, ToUserID: to_user_id}
	err = DB.Model(&model.Message{}).Create(message).Error
	return message, err
}

// 根据发送者ID和接收者ID获取消息列表
func FindMessagesBy_From_To_ID(ctx context.Context, from_user_id uint, to_user_id uint) (messages []model.Message, err error) {
	DB := GetDB(ctx)
	err = DB.Model(&model.Message{}).Where("from_user_id=? AND to_user_id=?", from_user_id, to_user_id).Preload(clause.Associations).Find(&messages).Error
	return messages, err
}
