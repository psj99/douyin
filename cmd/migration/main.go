package main

import (
	"context"
	"douyin/conf"
	"douyin/repository/dao"
	"douyin/utils"
)

func init() {
	conf.InitConfig()
	utils.InitLogger()
	dao.InitMySQL()
	//oss.InitOSS()
}

func main() {
	utils.PrintAsJson(conf.Cfg)
	utils.ZapLogger.Infof("asd")

	//初次使用或数据表结构变更时取消以下行的注释以迁移数据表
	m := NewMigrate(dao.GetDB(context.Background()))
	m.Run()

}
