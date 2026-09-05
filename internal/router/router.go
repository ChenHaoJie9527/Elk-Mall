package router

import (
	"github.com/ChenHaoJie9527/Elk-Mall/internal/controller"
	"github.com/labstack/echo/v5"
)

// Register 注册路由
func RegisterRouter(e *echo.Echo, h *controller.Health, u *controller.User) {
	// 健康检查
	e.GET("/ping", h.Ping)

	// 用户路由
	e.POST("/users/register", u.Register)
	e.POST("/users/login", u.Login)
	e.GET("/users/:id", u.GetByID)
}
