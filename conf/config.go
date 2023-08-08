package conf

import (
	"os"

	"github.com/spf13/viper"
)

var Cfg *Config

type Config struct {
	System *System `yaml:"system"`
	MySql  *MySql  `yaml:"mysql"`
	Redis  *Redis  `yaml:"redis"`
	Zap    *Zap    `yaml:"zap"`
	Log    *Log    `yaml:"log"`
}

type MySql struct {
	DbHost   string `yaml:"dbHost"`
	DbPort   string `yaml:"dbPort"`
	DbName   string `yaml:"dbName"`
	UserName string `yaml:"userName"`
	Password string `yaml:"password"`
	Charset  string `yaml:"charset"`
}

type Redis struct {
	RedisHost     string `yaml:"redisHost"`
	RedisPort     string `yaml:"redisPort"`
	RedisPassword string `yaml:"redisPassword"`
	RedisDbName   int    `yaml:"redisDbName"`
	RedisNetwork  string `yaml:"redisNetwork"`
}

type System struct {
	AppEnv   string `yaml:"appEnv"`
	Domain   string `yaml:"domain"`
	Version  string `yaml:"version"`
	HttpPort string `yaml:"httpPort"`
	Host     string `yaml:"host"`
}

type Zap struct {
	Level        string // 级别
	Prefix       string // 日志前缀
	Format       string // 输出
	Directory    string // 日志文件夹
	MaxAge       int    // 日志留存时间
	ShowLine     bool   // 显示行
	LogInConsole bool   // 输出控制台
}
type Log struct {
	MaxSize    int
	MaxBackups int
	MaxAge     int
	Compress   bool
}

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
