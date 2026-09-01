package controller

import (
	"net/http"

	"github.com/ChenHaoJie9527/Elk-Mall/internal/common/response"
	"github.com/labstack/echo/v5"
)

// Ping 健康检查接口
func Ping(c *echo.Context) error {
	return c.JSON(http.StatusOK, response.Success("pong"))
}
