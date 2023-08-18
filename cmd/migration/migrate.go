package main

import (
	"douyin/repository/model"
	"douyin/utils"
	"gorm.io/gorm"
)

type Migrate struct {
	db *gorm.DB
}

func NewMigrate(db *gorm.DB) *Migrate {
	return &Migrate{db: db}
}

func (m *Migrate) Run() {
	err := m.db.Set("gorm:table_options", "charset=utf8mb4").
		AutoMigrate(&model.User{}, &model.Video{}, &model.Comment{}, &model.Message{})
	if err != nil {
		panic("数据表迁移失败")
	} else {
		utils.ZapLogger.Info("AutoMigrate warn: 数据表迁移成功")
	}
}
