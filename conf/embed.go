package conf

import (
	"embed"
	_ "image"
)

// 读取assets下所有资源文件并嵌入程序
//
//go:embed assets
var Emb embed.FS
