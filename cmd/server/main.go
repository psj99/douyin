package main

import (
	"douyin/pkg/config"
	"douyin/pkg/http"
	"douyin/pkg/log"
	"go.uber.org/zap"
)

func main() {
	cfg := config.NewConfig()
	//printer.PrintAsJson(cfg)
	logger := log.NewLog(cfg)

	app, cleanup, err := NewApp(cfg, logger)
	if err != nil {
		panic(err)
	}

	logger.Info("server start",
		zap.String("host", "http://127.0.0.1"+cfg.System.HttpPort))

	//_ = app.Run(cfg.System.HttpPort)
	http.Run(app, cfg.System.HttpPort)
	defer cleanup()
}
