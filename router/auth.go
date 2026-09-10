package router

import (
	"context"
	"net/http"

	"github.com/ChenHaoJie9527/Elk-Mall/common"
	"github.com/ChenHaoJie9527/Elk-Mall/consts"
	"github.com/gin-gonic/gin"
)

type TokenFun func(ctx context.Context, token string) (*common.User, error)

func AuthMiddleware(filter func(*gin.Context) bool, getTokenFun TokenFun) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if filter != nil && !filter(ctx) {
			ctx.Next()
			return
		}

		token := ctx.GetHeader(consts.UserTokenKey)
		// 如果 token 为空，则返回 401 状态码
		if len(token) == 0 {
			// AbortWithStatusJSON 用于中断请求，并返回 JSON 响应，比单独使用 Abort() 更友好
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, common.AuthErr)
			return
		}

		// 获取用户信息
		user, err := getTokenFun(ctx, token)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, common.AuthErr.WithErr(err))
			return
		}

		// 设置用户信息到上下文
		ctx.Set(consts.CustomerUserKey, user)
		// 执行下一个中间件
		ctx.Next()
	}
}
