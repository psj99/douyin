package main

import (
	"douyin/conf"
	"douyin/repository/dao"
	"douyin/routes"
	"douyin/utils"
	"douyin/utils/oss"
)

func init() {
	conf.InitConfig()
	utils.InitLogger()
	dao.InitMySQL()
	oss.InitOSS()
}

func main() {

	r := routes.NewRouter()
	_ = r.Run(":" + conf.Cfg.System.HttpPort)
}
