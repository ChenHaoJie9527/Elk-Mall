package adaptor

import (
	"github.com/ChenHaoJie9527/Elk-Mall/config"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type IAdaptor interface {
	GetConfig() *config.Config
	GetMySQL() *gorm.DB
	GetRedis() *redis.Client
}

type Adaptor struct {
	conf  *config.Config
	mysql *gorm.DB
	redis *redis.Client
}

func NewAdaptor(conf *config.Config, mysql *gorm.DB, redis *redis.Client) *Adaptor {
	return &Adaptor{
		conf:  conf,
		mysql: mysql,
		redis: redis,
	}
}

func (a *Adaptor) GetConfig() *config.Config {
	return a.conf
}

func (a *Adaptor) GetMySQL() *gorm.DB {
	return a.mysql
}

func (a *Adaptor) GetRedis() *redis.Client {
	return a.redis
}
