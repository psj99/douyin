package oss

import (
	"douyin/conf"

	"context"
	"net"
	"net/url"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type MinIOService struct {
	client     *minio.Client
	bucketName string
}

func (m *MinIOService) init() {
	endpoint := net.JoinHostPort(conf.Cfg.OSS.OssHost, conf.Cfg.OSS.OssPort)
	accessKeyID := conf.Cfg.OSS.AccessKeyID
	secretAccessKey := conf.Cfg.OSS.SecretAccessKey

	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: false, // 设置是否使用TLS访问对象存储
	})
	if err != nil {
		panic(err)
	}

	m.client = client
	m.bucketName = conf.Cfg.OSS.BucketName
}

func (m *MinIOService) upload(ctx context.Context, objectName string, filePath string) (err error) {
	opts := minio.PutObjectOptions{} // 可选元数据
	_, err = m.client.FPutObject(ctx, m.bucketName, objectName, filePath, opts)
	return err
}

func (m *MinIOService) getURL(ctx context.Context, objectName string) (objectURL string, err error) {
	expires := time.Hour * time.Duration(conf.Cfg.OSS.Expiry).Abs()
	reqParams := make(url.Values) // 可选响应头
	presignedURL, err := m.client.PresignedGetObject(ctx, m.bucketName, objectName, expires, reqParams)
	return presignedURL.String(), err
}

func (m *MinIOService) remove(ctx context.Context, objectName string) (err error) {
	opts := minio.RemoveObjectOptions{} // 可选选项
	return m.client.RemoveObject(ctx, m.bucketName, objectName, opts)
}