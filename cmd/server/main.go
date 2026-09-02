package main

import (
	"fmt"
	"log"

	"github.com/ChenHaoJie9527/Elk-Mall/internal/common/response"
	"github.com/ChenHaoJie9527/Elk-Mall/internal/config"
	"github.com/ChenHaoJie9527/Elk-Mall/internal/router"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
)

func main() {
	cfg, err := config.Load("config.yaml")
	if err != nil {
		log.Fatalf("加载配置文件失败: %v", err)
	}

	fmt.Printf("配置文件加载成功: %+v\n", cfg)

	e := echo.New()
	// 挂载 自定义的全局错误处理函数
	e.HTTPErrorHandler = response.HTTPErrorHandler
	// 挂载 恢复中间件
	e.Use(middleware.Recover())
	// 挂载 路由
	router.Register(e)
	// 启动服务
	if err := e.Start(":" + cfg.Server.Port); err != nil {
		log.Fatal("start server: ", err)
	}
}
