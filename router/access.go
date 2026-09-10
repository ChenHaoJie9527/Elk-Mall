package router

import (
	"bytes"
	"io"
	"os"
	"time"

	"github.com/ChenHaoJie9527/Elk-Mall/consts"
	"github.com/ChenHaoJie9527/Elk-Mall/utils/logger"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// 访问记录过滤器
func (r *Router) AccessRecordFilter(ctx *gin.Context) bool {
	return true
}

// 获取请求体
func GetRequestBody(ctx *gin.Context) string {
	// TODO: 为什么要用 io.ReadAll 读取请求体？
	// TODO （答）: 因为 ctx.Request.Body 是一个 io.ReadCloser 接口，需要用 io.ReadAll 读取请求体
	body, _ := io.ReadAll(ctx.Request.Body)
	// 将请求体转换为字符串
	return string(body)
}

type responseWriterWrapper struct {
	gin.ResponseWriter           // 响应写入器
	Writer             io.Writer // 响应体缓冲区
}

type AccessLogFilter func(ctx *gin.Context) bool

// 访问日志中间件
func AccessLogMiddleware(filter AccessLogFilter) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 如果过滤器不为空，并且过滤器返回 false，则不记录访问日志
		if filter != nil && !filter(ctx) {
			ctx.Next()
			return
		}

		body := GetRequestBody(ctx)

		// 将请求体转换为 io.ReadCloser 接口 ，交给 请求体中间件 使用
		ctx.Request.Body = io.NopCloser(bytes.NewReader([]byte(body)))

		// 记录开始时间
		begin := time.Now()

		// 创建 zap 文件
		fields := CreateZapFiles(ctx, body)

		// 创建一个响应体缓冲区
		var responseBody bytes.Buffer

		// 创建一个多路复用器，将标准输出和响应体缓冲区合并
		multiWriter := io.MultiWriter(os.Stdout, &responseBody)

		// 创建一个响应体写入器
		ctx.Writer = &responseWriterWrapper{
			ResponseWriter: ctx.Writer,  // 响应写入器
			Writer:         multiWriter, // 响应体缓冲区
		}

		// 执行下一个中间件
		ctx.Next()

		// 获取响应体
		respBody := responseBody.String()

		// 如果响应体长度大于 1024，则截取前 1024 个字符
		if len(respBody) > 1024 {
			respBody = respBody[:1024]
		}

		fields = append(fields, zap.Int64("dur_ms", time.Since(begin).Milliseconds()))
		fields = append(fields, zap.Int("status", ctx.Writer.Status()))
		fields = append(fields, zap.String("resp", respBody))
		logger.Info("access_log", fields...)

	}
}

// 创建 zap 文件
func CreateZapFiles(ctx *gin.Context, body string) []zap.Field {
	fields := []zap.Field{
		zap.String("ip", ctx.RemoteIP()),
		zap.String("method", ctx.Request.Method),
		zap.String("path", ctx.Request.URL.Path),
		zap.String("params", ctx.Request.URL.RawQuery),
		zap.Any("body", body),
		zap.String("token", ctx.GetHeader(consts.UserTokenKey)),
	}

	return fields
}
