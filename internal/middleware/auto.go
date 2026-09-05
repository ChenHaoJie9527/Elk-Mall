package middleware

import (
	"github.com/ChenHaoJie9527/Elk-Mall/internal/pkg"
	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v5"
	"github.com/labstack/echo/v5"
)

// JWT 中间件: 解析 JWT 令牌
func JWT(secret []byte) echo.MiddlewareFunc {
	return echojwt.WithConfig(echojwt.Config{
		SigningKey:    secret,                 // 签名密钥
		SigningMethod: echojwt.AlgorithmHS256, // 签名方法
		NewClaimsFunc: func(c *echo.Context) jwt.Claims {
			return new(pkg.Claims) // 返回自定义的 Claims 指针类型
		},
	})
}
