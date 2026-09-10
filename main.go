package main

import (
	"context"
	"time"

	"github.com/ChenHaoJie9527/Elk-Mall/config"
	"github.com/ChenHaoJie9527/Elk-Mall/utils/logger"
	"github.com/redis/go-redis/v9"
	"github.com/samber/lo"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	conf := config.InitConfig()
	logger.SetLevel(conf.Server.LogLevel)

	_, err := initMySQL(&conf.MySQL)
	handleErr(err)
	logger.Debug("mysql connected successfully")

	_, err = initRedis(&conf.Redis)
	handleErr(err)
	logger.Debug("redis connected successfully")

	startServer(conf)
}

func startServer(conf *config.Config) {

}

// 初始化 MySQL 数据库连接
func initMySQL(conf *config.MySQL) (*gorm.DB, error) {
	// TODO: 为什么要设置最大空闲连接数和最大连接数？
	// 设置最大空闲连接数和最大连接数
	conf.MaxIdle = lo.Max([]int{conf.MaxIdle + 1, 5})
	conf.MaxOpen = lo.Max([]int{conf.MaxOpen + 1, 10})

	// 获取数据库连接字符串
	sqlDSN := conf.GetDSN()

	// 创建数据库连接
	db, err := gorm.Open(mysql.Open(sqlDSN))
	if err != nil {
		return nil, err
	}
	// 获取数据库连接池
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	// ping 数据库连接
	if err = sqlDB.Ping(); err != nil {
		return nil, err
	}

	sqlDB.SetMaxIdleConns(conf.MaxIdle)
	sqlDB.SetMaxOpenConns(conf.MaxOpen)

	return db, nil
}

func initRedis(conf *config.Redis) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:         conf.Addr,
		Password:     conf.PWD,
		DB:           conf.DBIndex,
		MaxIdleConns: conf.MaxIdle,
		PoolSize:     conf.MaxOpen,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	return client, nil
}

func handleErr(err error) {
	if err != nil {
		panic(err)
	}
}
