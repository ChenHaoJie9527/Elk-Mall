package config

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/ChenHaoJie9527/Elk-Mall/internal/logger"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	// 启用 Viper 的远程配置（etcd / consul 等），不 import 的话 ReadRemoteConfig 会直接报错
	// _ "github.com/spf13/viper/remote"
)

const (
	serverFullName = "elk.mall"
)

var (
	// etcd 里存放整份 YAML 的 key，和 v9 保持一致
	etcdKey = fmt.Sprintf("/config/%s/system", serverFullName)

	etcdAddr        string
	localConfigPath string
	source          string // 本次配置从哪来，仅用于启动日志
)

func init() {
	// -c 本地 YAML；-r / ETCD_ADDR 远程 etcd。有 etcd 地址时优先走远程。
	flag.StringVar(&localConfigPath, "c", "config.yaml", "本地 YAML 配置文件路径")
	flag.StringVar(&etcdAddr, "r", os.Getenv("ETCD_ADDR"), "etcd 地址，例如 http://127.0.0.1:2379")
}

// Config 配置
type Config struct {
	App    AppConfig   `mapstructure:"app"`
	Server AppServer   `mapstructure:"server"`
	MySQL  MySQLConfig `mapstructure:"mysql"`
	Redis  RedisConfig `mapstructure:"redis"`
	JWT    JWTConfig   `mapstructure:"jwt"`
}

type JWTConfig struct {
	Secret    string
	ExpiresIn time.Duration `mapstructure:"expires_in"`
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

type AppConfig struct {
	Name string
}

type AppServer struct {
	Port          string
	Env           string `mapstructure:"env"`
	LogLevel      string `mapstructure:"log_level"`
	LogFile       string `mapstructure:"log_file"`
	LogMaxSizeMB  int    `mapstructure:"log_max_size_mb"`
	LogMaxBackups int    `mapstructure:"log_max_backups"`
	LogMaxAgeDays int    `mapstructure:"log_max_age_days"`
}

// Source 返回本次加载来源，例如 "file:config.yaml" 或 "etcd:/config/elk.mall/system"
func Source() string {
	return source
}

// InitConfig 按优先级加载：etcd 地址非空 → 远程；否则 → 本地 YAML。
func InitConfig() (*Config, error) {
	// 如果 flag 未解析，则解析 flag
	if !flag.Parsed() {
		flag.Parse()
	}

	if strings.TrimSpace(etcdAddr) != "" {
		source = "etcd:" + etcdKey
		return loadFromEtcd()
	}

	source = "file:" + localConfigPath
	return loadFromFile(localConfigPath)
}

// Load 只读指定本地文件（给测试或明确只要文件的调用方用），不走 etcd。
func Load(path string) (*Config, error) {
	source = "file:" + path
	return loadFromFile(path)
}

// newViper 创建一个新的 Viper 实例
func newViper() *viper.Viper {
	vi := viper.New()
	vi.SetDefault("server.port", "8080")
	vi.SetDefault("server.env", "local")
	vi.SetDefault("server.log_level", "debug")
	vi.SetDefault("server.log_file", "./logs/server.log")
	vi.SetDefault("server.log_max_size_mb", 100)
	vi.SetDefault("server.log_max_backups", 3)
	vi.SetDefault("server.log_max_age_days", 28)
	vi.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	vi.AutomaticEnv()

	vi.SetDefault("mysql.max_open_conns", 10)
	vi.SetDefault("mysql.max_idle_conns", 5)

	vi.SetDefault("redis.addr", "127.0.0.1:6380")
	vi.SetDefault("redis.password", "")
	vi.SetDefault("redis.db", 0)

	vi.SetDefault("jwt.secret", "elk-mall-dev-jwt-secret")
	vi.SetDefault("jwt.expires_in", "1h")
	vi.SetDefault("server.log_level", "debug")
	return vi
}

func unmarshalConfig(vi *viper.Viper) (*Config, error) {
	var cfg Config
	if err := vi.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("解析配置失败: %w", err)
	}
	return &cfg, nil
}

// loadFromFile 读取本地 YAML 配置文件
func loadFromFile(path string) (*Config, error) {
	vi := newViper()
	vi.SetConfigFile(path)
	if err := vi.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("读取配置文件失败: %w", err)
	}
	return unmarshalConfig(vi)
}

// loadFromEtcd 读 etcd 的 YAML，并在后台按分钟轮询热更新。
// 热更新会改回这个指针里的字段；已经用旧值建好的对象（例如 JWT 中间件）不会自动重建。
func loadFromEtcd() (*Config, error) {
	endpoint := normalizeEtcdAddr(etcdAddr)

	vi := newViper()
	if err := vi.AddRemoteProvider("etcd3", endpoint, etcdKey); err != nil {
		return nil, fmt.Errorf("添加 etcd provider 失败: %w", err)
	}
	// etcd 的 value 没有文件名，必须显式告诉 Viper 是 YAML
	vi.SetConfigType("yaml")

	if err := vi.ReadRemoteConfig(); err != nil {
		return nil, fmt.Errorf("读取 etcd 配置失败 key=%s endpoint=%s: %w", etcdKey, endpoint, err)
	}

	cfg, err := unmarshalConfig(vi)
	if err != nil {
		return nil, err
	}

	go watchEtcd(vi, cfg)
	return cfg, nil
}

func watchEtcd(vi *viper.Viper, cfg *Config) {
	for {
		time.Sleep(time.Minute)
		if err := vi.WatchRemoteConfig(); err != nil {
			logger.Warn("etcd 配置热更新失败", zap.Error(err))
			continue
		}
		if err := vi.Unmarshal(cfg); err != nil {
			logger.Warn("etcd 配置反序列化失败", zap.Error(err))
			continue
		}
		logger.Info("etcd 配置已热更新",
			zap.String("app", cfg.App.Name),
			zap.String("port", cfg.Server.Port),
			zap.String("mysql_host", cfg.MySQL.Host),
			zap.String("redis_addr", cfg.Redis.Addr),
		)
	}
}

// Viper etcd provider 需要带 scheme 的地址，例如 http://127.0.0.1:2379
func normalizeEtcdAddr(addr string) string {
	addr = strings.TrimSpace(addr)
	if strings.Contains(addr, "://") {
		return addr
	}
	return "http://" + addr
}
