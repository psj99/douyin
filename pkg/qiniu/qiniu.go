package qiniu

import (
	"context"
	"douyin/conf"
	"encoding/base64"
	"fmt"
	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/storage"
	"log"
	"mime/multipart"
	"strconv"
	"strings"
	"time"
)

type QiniuUploader struct {
	AccessKey string
	SecretKey string
	Bucket    string
	Domain    string
}

func NewQiniuUploader(cfg *conf.Config) *QiniuUploader {
	return &QiniuUploader{
		AccessKey: cfg.Qiniu.AccessKey,
		SecretKey: cfg.Qiniu.SecretKey,
		Bucket:    cfg.Qiniu.Bucket,
		Domain:    cfg.Qiniu.Domain,
	}
}

func (uploader *QiniuUploader) UploadFile(ctx context.Context, file *multipart.FileHeader) (fileURL, coverURL string, err error) {
	//TODO 文件名先设置为了 时间戳 + filename, 后面应该还要与userid绑定
	file.Filename = strconv.FormatInt(time.Now().Unix(), 10) + file.Filename

	//图片路径
	filePreName := strings.Split(file.Filename, ".")[0]
	jpgPath := fmt.Sprintf("%s:%s.jpg", uploader.Bucket, filePreName)
	fmt.Println(jpgPath)
	saveJpgEntry := base64.URLEncoding.EncodeToString([]byte(jpgPath))
	vframeJpgFop := "vframe/jpg/offset/1|saveas/" + saveJpgEntry

	//连接多个操作指令
	persistentOps := strings.Join([]string{vframeJpgFop}, ";")
	pipeline := ""

	putPolicy := storage.PutPolicy{
		Scope:               uploader.Bucket,
		PersistentOps:       persistentOps,
		PersistentPipeline:  pipeline,
		PersistentNotifyURL: "http://api.example.com/qiniu/pfop/notify",
	}
	mac := qbox.NewMac(uploader.AccessKey, uploader.SecretKey)

	// 创建上传的凭证
	upToken := putPolicy.UploadToken(mac)

	// 打开上传文件
	src, err := file.Open()
	if err != nil {
		log.Println("打开上传文件发生了错误:    ", err)
		return
	}
	defer src.Close()

	// 上传文件到七牛云
	cfg := storage.Config{}
	formUploader := storage.NewFormUploader(&cfg)
	ret := storage.PutRet{}        // 上传后返回的结果
	putExtra := storage.PutExtra{} // 额外参数
	key := "" + file.Filename      // 上传路径，如果当前目录中已存在相同文件，则返回上传失败错误
	err = formUploader.Put(ctx, &ret, upToken, key, src, file.Size, &putExtra)

	if err != nil {
		log.Println("Put发生了错误:    ", err)
		return
	}
	fileURL = fmt.Sprintf("http://%s/%s.mp4", uploader.Domain, filePreName)

	coverURL = fmt.Sprintf("http://%s/%s.jpg", uploader.Domain, filePreName)
	return fileURL, coverURL, nil

}
