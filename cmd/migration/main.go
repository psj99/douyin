package main

import (
	"douyin/internal/repository"
	"douyin/pkg/config"
	"douyin/pkg/helper/printer"
	"douyin/pkg/log"
)

func main() {
	cfg := config.NewConfig()
	printer.PrintAsJson(cfg)

	m := NewMigrate(repository.NewDB(cfg), log.NewLog(cfg))
	m.Run()
}
