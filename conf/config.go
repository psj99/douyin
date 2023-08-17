package conf

import (
	"os"

	"github.com/spf13/viper"
)

type Config struct {
	System   *System   `yaml:"system"`
	MySQL    *MySQL    `yaml:"mysql"`
	Log      *Log      `yaml:"log"`
	MinioOSS *MinioOSS `yaml:"minio_oss" mapstructure:"minio_oss"`
	QiNiuOSS *QiNiuOSS `yaml:"qiniu_oss" mapstructure:"qiniu_oss"`
}

var Cfg *Config

func InitConfig() {
	workDir, _ := os.Getwd()
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(workDir + "/conf/locale")
	viper.AddConfigPath(workDir)
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	err = viper.Unmarshal(&Cfg)
	if err != nil {
		panic(err)
	}
}
