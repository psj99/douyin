package main

import (
	"douyin/conf"
	"douyin/internal/repository"
	"douyin/pkg/config"
	"douyin/pkg/helper/printer"
	"douyin/pkg/log"
)

var cfg *conf.Config

func init() {
	cfg = config.NewConfig()

}
func main() {
	printer.PrintAsJson(cfg)

	m := NewMigrate(repository.NewDB(cfg), log.NewLog(cfg))
	// 初次使用或数据表结构变更时取消以下行的注释以迁移数据表
	// dao.MakeMigrate()
	m.Run()
}
