package main

import (
	"github.com/ChenHaoJie9527/Elk-Mall/config"
	"github.com/ChenHaoJie9527/Elk-Mall/utils/logger"
)

func main() {
	conf := config.InitConfig()
	logger.SetLevel(conf.Server.LogLevel)
}
