package router

import (
	"github.com/ChenHaoJie9527/Elk-Mall/internal/controller"
	"github.com/labstack/echo/v5"
)

// Register 注册路由
func Register(e *echo.Echo) {
	e.GET("/ping", controller.Ping)
}
