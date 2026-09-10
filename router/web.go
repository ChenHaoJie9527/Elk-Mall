package router

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/ChenHaoJie9527/Elk-Mall/utils/logger"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type App struct {
	server *gin.Engine // 服务引擎
	addr   string      // 服务地址
}

func NewApp(port int, router IRouter) *App {

	gin.SetMode(gin.ReleaseMode) // 设置 Gin 模式为 ReleaseMode，这种模式 是生产环境，性能最佳，不会输出调试信息

	engine := gin.New() // 创建一个 Gin 引擎

	engine.Use(gin.Recovery()) // 使用 Gin 的 Recovery 中间件，用于恢复 panic 错误

	// 日志中间件,自定义过滤器，某些接口不需要记录日志
	engine.Use(AccessLogMiddleware(router.AccessRecordFilter))

	router.Register(engine)

	return &App{
		server: engine,
		addr:   ":" + strconv.Itoa(port),
	}
}

// 运行服务
func (a *App) Run() {
	service := &http.Server{
		Addr:    a.addr,
		Handler: a.server,
	}

	go func() {
		if err := service.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen err: %v", err)
		}
	}()

	logger.Debug(fmt.Sprintf("server started, endpoint: http://localhost%s", a.addr))

	// 创建一个通道，用于接收关闭信号
	closeCh := make(chan os.Signal, 1)
	// 接收关闭信号
	signal.Notify(closeCh, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	msg := <-closeCh

	logger.Warn("server closing: ", zap.String("msg", msg.String()))
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 关闭服务
	_ = service.Shutdown(ctx)
}
