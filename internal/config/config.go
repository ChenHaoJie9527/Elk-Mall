package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

// Config 配置
type Config struct {
	App    AppConfig
	Server AppServer
	MySQL  MySQLConfig
	Redis  RedisConfig
	JWT    JWTConfig
}

type JWTConfig struct {
	Secret string
}

type RedisConfig struct {
	Addr     string
	Password string
	DB       int
}

type MySQLConfig struct {
	Host         string
	Port         string
	User         string
	Password     string
	Database     string
	MaxOpenConns int `mapstructure:"max_open_conns"`
	MaxIdleConns int `mapstructure:"max_idle_conns"`
}

// AppConfig 应用配置
type AppConfig struct {
	Name string
}

// AppServer 应用服务器配置
type AppServer struct {
	Port string
}

// Load 加载配置文件
func Load(path string) (*Config, error) {
	// 初始化 viper
	vi := viper.New()
	// 设置配置文件的路径
	vi.SetConfigFile(path)
	// 设置默认值，如果配置文件中没有设置，则使用默认值
	vi.SetDefault("server.port", "8080")
	// 将 server.port 转换为 SERVER_PORT，用于环境变量
	vi.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	// 自动加载环境变量
	vi.AutomaticEnv()

	// 设置默认值，如果配置文件中没有设置，则使用默认值
	vi.SetDefault("mysql.max_open_conns", 10)
	vi.SetDefault("mysql.max_idle_conns", 5)

	vi.SetDefault("redis.addr", "127.0.0.1:6380")
	vi.SetDefault("redis.password", "")
	vi.SetDefault("redis.db", 0)

	// JWT 配置
	vi.SetDefault("jwt.secret", "elk-mall-dev-jwt-secret")

	// 读取配置文件
	if err := vi.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("读取配置文件失败: %v", err)
	}

	var cfg Config
	// 将配置文件中的配置解析到 cfg 中
	if err := vi.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("解析配置文件失败: %v", err)
	}

	return &cfg, nil
}
