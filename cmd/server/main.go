package main

import (
	"fmt"
	"log"

	"github.com/ChenHaoJie9527/Elk-Mall/internal/config"
	"github.com/ChenHaoJie9527/Elk-Mall/internal/router"
	"github.com/labstack/echo/v5"
)

func main() {
	cfg, err := config.Load("config.yaml")
	if err != nil {
		log.Fatalf("加载配置文件失败: %v", err)
	}

	fmt.Printf("配置文件加载成功: %+v\n", cfg)

	e := echo.New()
	router.Register(e)

	if err := e.Start(":" + cfg.Server.Port); err != nil {
		log.Fatal("start server: ", err)
	}

}
