package main

import (
	"os"

	"github.com/ChenHaoJie9527/Elk-Mall/internal/logger"
	"github.com/ChenHaoJie9527/Elk-Mall/internal/server"
	"go.uber.org/zap"
)

func main() {
	// 启动逻辑在 server.Run：失败时这里统一打日志并退出。
	// 好处是 main 里不再到处 log.Fatal，defer Close() 也能在 Run 返回前执行到。
	if err := server.Run(); err != nil {
		logger.Error("进程退出", zap.Error(err))
		logger.Sync()
		os.Exit(1)
	}
}
