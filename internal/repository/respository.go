package repository

import (
	"context"
	"douyin/conf"
	"fmt"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

type Repository struct {
	db     *gorm.DB
	rdb    *redis.Client
	logger *log.Logger
}

func NewRepository(db *gorm.DB, rdb *redis.Client, logger *log.Logger) *Repository {
	return &Repository{
		db:     db,
		rdb:    rdb,
		logger: logger,
	}
}

func NewDB(cfg *conf.Config) *gorm.DB {
	mysqlCfg := cfg.MySql

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=True&loc=Local",
		mysqlCfg.UserName,
		mysqlCfg.Password,
		mysqlCfg.DbHost,
		mysqlCfg.DbPort,
		mysqlCfg.DbName,
		mysqlCfg.Charset,
	)

	var ormLogger = logger.Default
	if gin.Mode() == "debug" {
		ormLogger = logger.Default.LogMode(logger.Info)
	}
	// db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN:                       dsn,   // DSN data source name
		DefaultStringSize:         256,   // string 类型字段的默认长度
		DisableDatetimePrecision:  true,  // 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
		DontSupportRenameIndex:    true,  // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
		DontSupportRenameColumn:   true,  // 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
		SkipInitializeWithVersion: false, // 根据版本自动配置
	}), &gorm.Config{
		Logger: ormLogger, // 打印日志
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, // 表明不加s
		},
	})

	if err != nil {
		panic(err)
	}

	sqlDB, _ := db.DB()
	sqlDB.SetMaxIdleConns(mysqlCfg.MaxIdleConns)                                    // 设置连接池中的最大闲置连接
	sqlDB.SetMaxOpenConns(mysqlCfg.MaxOpenConns)                                    // 设置数据库的最大连接数量
	sqlDB.SetConnMaxLifetime(time.Duration(mysqlCfg.ConnMaxLifetime) * time.Second) // 设置连接的最大可复用时间
	return db
}

func NewRedis(cfg *conf.Config) *redis.Client {
	redisCfg := cfg.Redis
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", redisCfg.RedisHost, redisCfg.RedisPort),
		Password: redisCfg.RedisPassword,
		DB:       redisCfg.RedisDbName,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		panic(fmt.Sprintf("redis error: %s", err.Error()))
	}

	return rdb
}
