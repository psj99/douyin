package service

import (
	"douyin/pkg/helper/qiniu"
	"douyin/pkg/helper/sid"
	"douyin/pkg/jwt"
	"douyin/pkg/log"
)

type Service struct {
	logger   *log.Logger
	sid      *sid.Sid
	jwt      *jwt.JWT
	uploader *qiniu.QiniuUploader
}

func NewService(logger *log.Logger, sid *sid.Sid, jwt *jwt.JWT, uploader *qiniu.QiniuUploader) *Service {
	return &Service{
		logger:   logger,
		sid:      sid,
		jwt:      jwt,
		uploader: uploader,
	}
}
