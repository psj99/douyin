package dao

import (
	"douyin/repository/model"
	"log"

	"gorm.io/gorm"
)

// MakeMigrate 迁移数据表，在没有数据表结构变更时候，建议注释不执行
// 只支持创建表、增加表中没有的字段和索引
// 为了保护数据，并不支持改变已有的字段类型或删除未被使用的字段
func MakeMigrate(db *gorm.DB) {
	err := db.Set("gorm:table_options", "charset=utf8mb4").
		AutoMigrate(&model.User{}, &model.Video{}, &model.Comment{})

	if err != nil {
		panic("gorm 迁移失败")
	} else {
		log.Println("gorm 迁移成功")
	}

}
