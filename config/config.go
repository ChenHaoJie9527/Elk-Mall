package config

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/gogf/gf/util/gconv"
	"github.com/spf13/viper"
	_ "github.com/spf13/viper/remote"
	"go.yaml.in/yaml/v3"
)

const (
	ServerName     = "mall"
	FullServerName = "elk.mall"
)

var (
	etcdKey         = fmt.Sprintf("/configs/%s/system", FullServerName)
	etcdAddr        string // etcd 地址
	localConfigPath string // 本地配置配置文件路径
	GlobalConfig    Config // 全局配置
)

type Config struct {
	Server Server `yaml:"server"`
	MySQL  MySQL  `yaml:"mysql"`
	Redis  Redis  `yaml:"redis"`
}

type Server struct {
	HttpPort    int    `yaml:"http_port"`
	Env         string `yaml:"env"`
	EnablePprof bool   `yaml:"enable_pprof"`
	LogLevel    string `yaml:"log_level"`
}

type MySQL struct {
	Dialect  string `yaml:"dialect"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Database string `yaml:"database"`
	Charset  string `yaml:"charset"`
	ShowSql  bool   `yaml:"show_sql"`
	MaxOpen  int    `yaml:"max_open"`
	MaxIdle  int    `yaml:"max_idle"`
}

// GetDSN 获取 MySQL 的 DSN -
// 格式：user:password@tcp(host:port)/database?charset=charset&parseTime=true&loc=Local
func (sql *MySQL) GetDSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=true&loc=Local",
		sql.User, sql.Password, sql.Host, sql.Port, sql.Database, sql.Charset)
}

type Redis struct {
	Addr    string `yaml:"addr"`
	PWD     string `yaml:"password"`
	DBIndex int    `yaml:"db_index"`
	MaxIdle int    `yaml:"max_idle"`
	MaxOpen int    `yaml:"max_open"`
}

// 注册 -c（本地 yaml）和 -r（etcd 地址，默认读 ETCD_ADDR）
func init() {
	flag.StringVar(&localConfigPath, "c", ServerName+"_config.yaml", "local yaml config path")
	flag.StringVar(&etcdAddr, "r", os.Getenv("ETCD_ADDR"), "etcd endpoint, e.g. http://127.0.0.1:2379")
}

func InitConfig() *Config {

	var (
		err      error         // 错误变量
		tempConf = &Config{}   // 配置变量
		vipConf  = viper.New() // viper 配置变量
	)

	// 解析命令行参数
	if !flag.Parsed() {
		flag.Parse()
	}

	if etcdAddr != "" {
		tempConf, err = getFromRemoteAndWatchUpdate(vipConf)
		if err != nil {
			panic(err)
		}
		return tempConf
	}

	// 从本地获取
	tempConf, err = getFromLocal()
	if err != nil {
		panic(err)
	}
	return tempConf
}

// 当 存在 etcd时，从 etcd 读取配置，并监听配置变化
func getFromRemoteAndWatchUpdate(v *viper.Viper) (*Config, error) {
	tempConf := Config{}
	if err := prepareRemoteConfig(v, etcdAddr); err != nil {
		return nil, err
	}

	if err := v.ReadRemoteConfig(); err != nil {
		return nil, fmt.Errorf("read etcd config (endpoint=%s, key=%s): %w", etcdAddr, etcdKey, err)
	}

	if err := v.Unmarshal(&tempConf); err != nil {
		return nil, err
	}

	go func() {
		for {
			time.Sleep(time.Minute)
			if err := v.WatchRemoteConfig(); err != nil {
				continue
			}
			if err := v.Unmarshal(&GlobalConfig); err != nil {
				continue
			}
			fmt.Println(">>> etcd config hot-reloaded: ", gconv.String(GlobalConfig))
		}
	}()

	return &tempConf, nil
}

// AddRemoteProvider 参数顺序是 provider, endpoint, path。
// etcd 的 endpoint 必须是 http://ip:port，不能把配置 key 或本地 yaml 路径当成地址。
func prepareRemoteConfig(v *viper.Viper, endpoint string) error {
	if looksLikeConfigFile(endpoint) {
		return fmt.Errorf("-r 需要 etcd 地址（例如 http://127.0.0.1:2379），本地配置请使用 -c %s", endpoint)
	}
	v.SetConfigType("yaml")
	return v.AddRemoteProvider("etcd3", normalizeEtcdEndpoint(endpoint), etcdKey)
}

func looksLikeConfigFile(s string) bool {
	lower := strings.ToLower(s)
	return strings.HasSuffix(lower, ".yaml") || strings.HasSuffix(lower, ".yml") ||
		strings.HasSuffix(lower, ".json") || strings.HasSuffix(lower, ".toml")
}

func normalizeEtcdEndpoint(addr string) string {
	if strings.Contains(addr, "://") {
		return addr
	}
	return "http://" + addr
}

// 当 不存在 etcd时，从本地文件读取配置
func getFromLocal() (*Config, error) {
	tempConf := Config{}
	content, err := os.ReadFile(localConfigPath)
	if err != nil {
		return nil, fmt.Errorf("local config file not found: %s", localConfigPath)
	}
	if err = yaml.Unmarshal(content, &tempConf); err != nil {
		return nil, err
	}
	return &tempConf, nil
}
