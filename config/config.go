package config

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/gogf/gf/util/gconv"
	"github.com/spf13/viper"
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
	http_port   int    `yaml:"http_port"`
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

// 注册 -c（默认 mall_local.yml）和 -r（etcd，默认读 ETCD_ADDR）
func init() {
	flag.StringVar(&localConfigPath, "c", ServerName+"_config.yaml", "default config path")
	flag.StringVar(&etcdAddr, "r", os.Getenv("ETCD_ADDR"), "default etcd address")
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
	// 添加 etcd 远程提供者
	if err := v.AddRemoteProvider("etcd3", etcdKey, etcdAddr); err != nil {
		return nil, err
	}

	// 读取配置
	if err := v.ReadRemoteConfig(); err != nil {
		return nil, err
	}

	// 将配置反序列化到 tempConf
	if err := v.Unmarshal(&tempConf); err != nil {
		return nil, err
	}

	// 监听配置变化
	go func() {
		for {
			// 每隔 1 分钟 监听一次配置变化
			time.Sleep(time.Minute * 1)
			if err := v.WatchRemoteConfig(); err != nil {
				v.Unmarshal(&GlobalConfig)
				// 将配置转换为字符串并打印
				fmt.Println(">>> etcd config hot-reloaded: ", gconv.String(GlobalConfig))
			}
		}
	}()

	return &tempConf, nil

}

// 当 不存在 etcd时，从本地文件读取配置
func getFromLocal() (*Config, error) {
	tempConf := Config{}
	// 如果本地配置文件存在，则读取配置文件
	if _, err := os.Stat(localConfigPath); err == nil {
		content, err := os.ReadFile(localConfigPath)
		if err != nil {
			return nil, err
		}
		err = yaml.Unmarshal(content, &tempConf)
		if err != nil {
			return nil, err
		}
	}
	return &tempConf, nil
}
