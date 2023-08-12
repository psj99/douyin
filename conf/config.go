package conf

import (
	"flag"
	"fmt"
	"os"

	"github.com/spf13/viper"
)

type System struct {
	AppEnv   string `yaml:"appEnv"`
	HttpPort string `yaml:"httpPort"`
}

type JWT struct {
	SignKey string `yaml:"sign_key"`
}

type MySql struct {
	DbHost          string `yaml:"dbHost"`
	DbPort          string `yaml:"dbPort"`
	DbName          string `yaml:"dbName"`
	UserName        string `yaml:"userName"`
	Password        string `yaml:"password"`
	Charset         string `yaml:"charset"`
	MaxIdleConns    int    `yaml:"maxIdleConns"`
	MaxOpenConns    int    `yaml:"maxOpenConns"`
	ConnMaxLifetime int    `yaml:"connMaxLifetime"`
}

type Redis struct {
	RedisHost     string `yaml:"redisHost"`
	RedisPort     string `yaml:"redisPort"`
	RedisPassword string `yaml:"redisPassword"`
	RedisDbName   int    `yaml:"redisDbName"`
	RedisNetwork  string `yaml:"redisNetwork"`
}
type Log struct {
	Level        string `yaml:"level"`
	Prefix       string `yaml:"prefix"`
	Format       string `yaml:"format"`
	Directory    string `yaml:"directory"`
	MaxAge       int    `yaml:"maxAge"`
	ShowLine     bool   `yaml:"showLine"`
	LogInConsole bool   `yaml:"logInConsole"`
	MaxSize      int    `yaml:"maxSize"`
	MaxBackups   int    `yaml:"maxBackups"`
	Compress     bool   `yaml:"compress"`
}

type Config struct {
	System *System `yaml:"system"`
	MySql  *MySql  `yaml:"mysql"`
	Redis  *Redis  `yaml:"redis"`
	Log    *Log    `yaml:"log"`
}

var Cfg *Config

func InitConfig() {

	viper := NewConfig()
	err := viper.Unmarshal(&Cfg)
	if err != nil {
		panic(err)
	}
}

func NewConfig() *viper.Viper {
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
