package main

import (
	"douyin/conf"
	"douyin/repository/dao"
	"douyin/routes"
	"douyin/utils"
)

func init() {
	conf.InitConfig()
	utils.InitLogger()
	dao.MySQLInit()
}

func main() {
	utils.PrintAsJson(conf.Cfg)
	r := routes.NewRouter()
	_ = r.Run(conf.Cfg.System.HttpPort)
}
