package router

import (
	"github.com/ChenHaoJie9527/Elk-Mall/internal/controller"
	"github.com/labstack/echo/v5"
)

// Register 注册路由
// e: echo 实例
// h: 健康检查控制器
// u: 用户控制器
// jwtMW: JWT 中间件
func RegisterRouter(e *echo.Echo, h *controller.Health, u *controller.User, jwtMW echo.MiddlewareFunc) {
	// 健康检查
	e.GET("/ping", h.Ping)

	// 公开的路由，不需要认证
	e.POST("/users/register", u.Register)
	e.POST("/users/login", u.Login)

	// 需要认证的路由组，使用 JWT 中间件
	auth := e.Group("", jwtMW)
	// 获取用户信息 需要认证
	auth.GET("/users/:id", u.GetByID)
}
