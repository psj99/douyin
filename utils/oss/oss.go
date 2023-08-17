package oss

import (
	"douyin/conf"
	"douyin/storage"
	"douyin/utils"

	"context"
	"errors"
	"io"
	"os"
	"path/filepath"
)

// 自定义错误类型
var ErrorRollbackFailed = errors.New("回滚操作(封面移除)失败")

// 自定义对象扩展名 需包含"."
const videoExt = ".mp4"
const coverExt = ".png"

// 欢迎添加对其他OSS的支持
type OSService interface {
	init()
	upload(ctx context.Context, objectName string, filePath string) (err error)  // 上传对象
	getURL(ctx context.Context, objectName string) (objectURL string, err error) // 获取对象外链
	remove(ctx context.Context, objectName string) (err error)                   // 移除对象

	// 以下为流式上传方案所需
	uploadStream(ctx context.Context, objectName string, reader io.Reader, objectSize int64) (err error) // 流式上传对象
	download(ctx context.Context, objectName string, filePath string) (err error)                        // 下载对象
}

var _oss OSService

func InitOSS() {
	if conf.Cfg.MinioOSS.Service == "minio" {
		_oss = &MinIOService{}
	} else {
		panic(errors.New("暂不支持该OSS: " + conf.Cfg.MinioOSS.Service))
	}
	_oss.init()
}

// 获取对象在存储桶内所用名称
func GetObjectName(objectID string) (videoName string, coverName string) {
	return objectID + videoExt, objectID + coverExt
}

// 将filePath处的视频上传为唯一标识为objectID的对象 自动切取封面并以并以同名上传
func UploadVideo(ctx context.Context, objectID string, videoPath string) (err error) {
	// 视频对象与封面对象名
	videoName, coverName := GetObjectName(objectID)

	// 切取封面
	coverPath := filepath.Join(conf.Cfg.System.TempDir, coverName) // 临时文件位置
	err = utils.GetSnapshot(videoPath, coverPath, 1)               // 防止切取黑屏
	if err != nil {
		utils.ZapLogger.Errorf("GetSnapshot err: %v", err)
		return err
	}
	defer os.Remove(coverPath) // 不保证自动清理成功 但临时数据在本地 易于检测是否仍存在且可被直接复写

	// 上传
	err = _oss.upload(ctx, coverName, coverPath) // 先传输小文件
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

// 以下为流式上传方案所需
// 流式上传(视频)对象
func UploadVideoStream(ctx context.Context, objectID string, videoStream io.Reader, videoSize int64) (err error) {
	// 视频对象与封面对象名
	videoName, coverName := GetObjectName(objectID)

	// 获取默认封面
	coverStream, err := storage.Emb.Open("assets/defaultCover" + coverExt)
	if err != nil {
		utils.ZapLogger.Errorf("Emb.Open (defaultCover) err: %v", err)
		return err
	}
	defer coverStream.Close() // 不保证自动关闭成功

	coverStat, err := coverStream.Stat()
	if err != nil {
		utils.ZapLogger.Errorf("File.Stat (defaultCover) err: %v", err)
		return err
	}
	coverSize := coverStat.Size()

	// 上传
	err = _oss.uploadStream(ctx, coverName, coverStream, coverSize) // 先传输小文件
	if err != nil {
		utils.ZapLogger.Errorf("_oss.uploadStream (cover) err: %v", err)
		return err
	}
	err = _oss.uploadStream(ctx, videoName, videoStream, videoSize) // 视频传输失败时将移除其封面
	if err != nil {
		utils.ZapLogger.Errorf("_oss.uploadStream (video) err: %v", err)
		utils.ZapLogger.Warnf("_oss.uploadStream (cover) warn: 正在回滚(移除对应封面%v)", coverName)
		err2 := _oss.remove(ctx, coverName)
		if err2 != nil {
			return ErrorRollbackFailed
		} else {
			return err
		}
	}
	return err
}

// 更新封面
func UpdateCover(ctx context.Context, objectID string) (err error) {
	// 视频对象与封面对象名
	videoName, coverName := GetObjectName(objectID)

	// 下载视频对象到本地
	videoPath := filepath.Join(conf.Cfg.System.TempDir, videoName)
	err = _oss.download(ctx, videoName, videoPath)
	if err != nil {
		utils.ZapLogger.Errorf("_oss.download (video) err: %v", err)
		return err
	}
	defer os.Remove(videoPath) // 不保证自动清理成功 但临时数据在本地 易于检测是否仍存在且可被直接复写

	// 切取封面
	coverPath := filepath.Join(conf.Cfg.System.TempDir, coverName) // 临时文件位置
	err = utils.GetSnapshot(videoPath, coverPath, 1)               // 防止切取黑屏
	if err != nil {
		utils.ZapLogger.Errorf("GetSnapshot err: %v", err)
		return err
	}
	defer os.Remove(coverPath) // 不保证自动清理成功 但临时数据在本地 易于检测是否仍存在且可被直接复写

	// 上传
	err = _oss.upload(ctx, coverName, coverPath)
	if err != nil {
		utils.ZapLogger.Errorf("_oss.upload (cover) err: %v", err)
		return err
	}
	utils.ZapLogger.Infof("UpdateCover info: %v - success", coverName)
	return err
}
