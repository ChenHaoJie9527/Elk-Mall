package controller

import (
	"database/sql"
	"net/http"

	"github.com/ChenHaoJie9527/Elk-Mall/internal/common/response"
	"github.com/labstack/echo/v5"
	"github.com/redis/go-redis/v9"
)

type Health struct {
	MySQL *sql.DB
	Redis *redis.Client
}

// Ping 健康检查接口
func (h *Health) Ping(c *echo.Context) error {
	// 把请求上下文传递给 MySQL 连接，确保在请求上下文被取消时，MySQL 连接也会被关闭
	if err := h.MySQL.PingContext(c.Request().Context()); err != nil {
		return err
	}
	// 把请求上下文传递给 Redis 连接，确保在请求上下文被取消时，Redis 连接也会被关闭
	if err := h.Redis.Ping(c.Request().Context()).Err(); err != nil {
		return err
	}
	// 返回成功响应
	return c.JSON(http.StatusOK, response.Success(map[string]string{
		"mysql": "ok",
		"redis": "ok",
	}))
}
