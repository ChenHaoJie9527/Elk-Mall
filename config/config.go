package config

import (
	"flag"
	"fmt"
	"os"
)

const (
	ServerName     = "mall"
	FullServerName = "elk.mall"
)

var (
	etcd            = fmt.Sprintf("/configs/%s/system", FullServerName)
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

// init 初始化配置
func init() {
	// 解析命令行参数 格式：-config=config.yaml
	flag.StringVar(&localConfigPath, "config", ServerName+"_config.yaml", "default config path")
	// 解析命令行参数 格式：-r=127.0.0.1:2379
	flag.StringVar(&etcdAddr, "r", os.Getenv("ETCD_ADDR"), "default etcd address")
}
