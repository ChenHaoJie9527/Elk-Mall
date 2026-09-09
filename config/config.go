package config

import "fmt"

const (
	ServerName     = "Mall"
	FullServerName = "elk.mall"
)

var (
	etcd            = fmt.Sprintf("/configs/%s/system", FullServerName)
	etcdKey         string // etcd 配置的 key
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

type Redis struct {
	Addr    string `yaml:"addr"`
	PWD     string `yaml:"password"`
	DBIndex int    `yaml:"db_index"`
	MaxIdle int    `yaml:"max_idle"`
	MaxOpen int    `yaml:"max_open"`
}
