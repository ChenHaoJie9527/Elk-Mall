package router

import (
	"github.com/ChenHaoJie9527/Elk-Mall/internal/controller"
	"github.com/labstack/echo/v5"
)

func RegisterRouter(e *echo.Echo, h *controller.Health, u *controller.User, jwtMW echo.MiddlewareFunc) {
	e.GET("/ping", h.Ping)

	e.POST("/users/register", u.Register)
	e.POST("/users/login", u.Login)

	auth := e.Group("", jwtMW)
	auth.GET("/users/me", u.Me)
	auth.GET("/users/:id", u.GetByID)
}
