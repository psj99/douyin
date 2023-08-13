package main

import (
	"douyin/internal/repository"
	"douyin/pkg/config"
	"douyin/pkg/helper/printer"
	"douyin/pkg/log"
)

//var cfg *conf.Config

func main() {
	cfg := config.NewConfig()
	printer.PrintAsJson(cfg)

	m := NewMigrate(repository.NewDB(cfg), log.NewLog(cfg))
	// 初次使用或数据表结构变更时取消以下行的注释以迁移数据表
	// dao.MakeMigrate()
	m.Run()
}
