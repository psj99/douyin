package oss

import (
	"douyin/conf"
	"douyin/utils"

	"context"
	"errors"
	"path/filepath"
)

// 自定义错误类型
var ErrorRollbackFailed = errors.New("回滚操作(封面移除)失败")

// 欢迎添加对其他OSS的支持
type OSService interface {
	init()
	upload(ctx context.Context, objectName string, filePath string) (err error)
	getURL(ctx context.Context, objectName string) (objectURL string, err error)
	remove(ctx context.Context, objectName string) (err error)
}

var _oss OSService

func InitOSS() {
	if conf.Cfg.OSS.Service == "minio" {
		_oss = &MinIOService{}
	} else {
		panic(errors.New("暂不支持该OSS: " + conf.Cfg.OSS.Service))
	}
	_oss.init()
}

// 将filePath处的视频上传为唯一标识为objectID的对象 自动获取封面并以并以同名上传
func UploadVideo(ctx context.Context, objectID string, videoPath string) (err error) {
	// 切取封面
	snapshotPath := videoPath[:len(videoPath)-len(filepath.Ext(videoPath))] + ".png"
	utils.GetSnapshot(videoPath, snapshotPath, 24) // 取24帧格式第二秒第一帧防止切取黑屏

	// 视频对象与封面对象名
	videoName, coverName := GetObjectName(objectID)

	// 上传
	err = _oss.upload(ctx, coverName, snapshotPath) // 先传输小文件
	if err != nil {
		utils.ZapLogger.Errorf("_oss.upload (cover) err: %v", err)
		return err
	}
	err = _oss.upload(ctx, videoName, videoPath) // 视频传输失败时将移除其封面
	if err != nil {
		utils.ZapLogger.Errorf("_oss.upload (video) err: %v", err)
		utils.ZapLogger.Warnf("_oss.upload (cover) warn: 正在回滚(移除对应封面%v)", coverName)
		err2 := _oss.remove(ctx, coverName)
		if err2 != nil {
			return ErrorRollbackFailed
		} else {
			return err
		}
	}
	return err
}

// 获取唯一标识为objectID的(视频)对象的短期外链 自动获取同名封面的短期外链
func GetVideo(ctx context.Context, objectID string) (videoURL string, coverURL string, err error) {
	// 视频对象与封面对象名
	videoName, coverName := GetObjectName(objectID)

	// 获取URL
	videoURL, err = _oss.getURL(ctx, videoName)
	if err != nil {
		utils.ZapLogger.Errorf("_oss.getURL (video) err: %v", err)
		return "", "", err
	}
	coverURL, err = _oss.getURL(ctx, coverName)
	if err != nil {
		utils.ZapLogger.Errorf("_oss.getURL (cover) err: %v", err)
		return videoURL, "", err
	}
	return videoURL, coverURL, err
}

// 移除指定对象(需指定完整对象名) 可用于手动错误回滚
func RemoveObject(ctx context.Context, objectName string) (err error) {
	return _oss.remove(ctx, objectName)
}

// 获取对象在存储桶内所用名称 可用于手动错误回滚
func GetObjectName(objectID string) (videoName string, coverName string) {
	return objectID + ".mp4", objectID + ".png"
}
