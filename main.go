package main

import (
	"douyin/conf"
	"douyin/pkg/utils"
	"douyin/repository/dao"
	"douyin/routes"
)

func init() {
	conf.InitConfig()
	utils.InitLogger()
	dao.MySQLInit()
}

func main() {

	utils.PrintAsJson(conf.Cfg)
	utils.ZapLogger.Info("nihao")
	//db := dao.NewDBClient(context.Background())
	//dao.MakeMigrate(db)
	r := routes.NewRouter()
	_ = r.Run(conf.Cfg.System.HttpPort)
}
