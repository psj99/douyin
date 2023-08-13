//go:build wireinject
// +build wireinject

package main

import (
	"douyin/conf"
	"douyin/internal/handler"
	"douyin/internal/repository"
	"douyin/internal/router"
	"douyin/internal/service"
	"douyin/pkg/helper/sid"
	"douyin/pkg/jwt"
	"douyin/pkg/log"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

var HandlerSet = wire.NewSet(
	handler.NewHandler,
	handler.NewUserHandler,
)

var ServiceSet = wire.NewSet(
	service.NewService,
	service.NewUserService,
)

var RepositorySet = wire.NewSet(
	repository.NewDB,
	repository.NewRedis,
	repository.NewRepository,
	repository.NewUserRepository,
)

func NewApp(*conf.Config, *log.Logger) (*gin.Engine, func(), error) {
	panic(wire.Build(
		RepositorySet,
		ServiceSet,
		HandlerSet,
		router.NewRouter,
		sid.NewSid,
		jwt.NewJwt,
	))
}
