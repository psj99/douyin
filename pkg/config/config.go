package config

import (
	"douyin/conf"
	"flag"
	"fmt"
	"os"

	"github.com/spf13/viper"
)

func NewConfig() *conf.Config {
	cfg := &conf.Config{}
	viper := viperLoad()
	err := viper.Unmarshal(cfg)
	if err != nil {
		panic(err)
	}
	return cfg
}

func viperLoad() *viper.Viper {
	// 配置文件优先级：环境变量 > 命令行参数 > 默认值
	envConf := os.Getenv("APP_CONF")
	if envConf == "" {
		flag.StringVar(&envConf, "conf", "conf/locale/config.yaml", "config path, eg: -conf conf/local.yaml")
		flag.Parse()
	}
	fmt.Println("load conf file:", envConf)

	v := viper.New()
	v.SetConfigFile(envConf)
	err := v.ReadInConfig()
	if err != nil {
		panic(err)
	}

	return v
}
