package controller

import (
	"database/sql"
	"net/http"

	"github.com/ChenHaoJie9527/Elk-Mall/internal/common/response"
	"github.com/labstack/echo/v5"
)

type Health struct {
	MySQL *sql.DB
}

// Ping 健康检查接口
func (h *Health) Ping(c *echo.Context) error {
	if err := h.MySQL.PingContext(c.Request().Context()); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, response.Success(map[string]string{
		"mysql": "ok",
	}))
}
