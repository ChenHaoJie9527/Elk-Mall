package main

import (
	"fmt"
	"log"
	"log/slog"
	"os"

	"github.com/ChenHaoJie9527/Elk-Mall/internal/adaptor"
	"github.com/ChenHaoJie9527/Elk-Mall/internal/common/response"
	"github.com/ChenHaoJie9527/Elk-Mall/internal/config"
	"github.com/ChenHaoJie9527/Elk-Mall/internal/controller"
	authmw "github.com/ChenHaoJie9527/Elk-Mall/internal/middleware"
	"github.com/ChenHaoJie9527/Elk-Mall/internal/model/do"
	"github.com/ChenHaoJie9527/Elk-Mall/internal/repository"
	"github.com/ChenHaoJie9527/Elk-Mall/internal/router"
	"github.com/ChenHaoJie9527/Elk-Mall/internal/service"
	"github.com/google/uuid"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
)

func main() {
	cfg, err := config.Load("config.yaml")
	if err != nil {
		log.Fatalf("加载配置文件失败: %v", err)
	}

	fmt.Printf("配置文件加载成功: %+v\n", cfg)

	gdb, err := adaptor.OpenMySql(cfg.MySQL)
	if err != nil {
		log.Fatalf("连接 MySQL 失败: %v", err)
	}

	// 从 GORM 拿回原来的 *sql.DB
	sqlDB, err := gdb.DB()
	if err != nil {
		log.Fatalf("取出 MySQL 连接失败: %v", err)
	}
	// Close 在 *sql.DB 上，*gorm.DB 没有这个方法
	defer sqlDB.Close()

	// 按 UserDO 建/补 users 表
	if err := gdb.AutoMigrate(&do.UserDO{}); err != nil {
		log.Fatalf("自动迁移失败: %v", err)
	}

	r, err := adaptor.OpenRedis(cfg.Redis)
	if err != nil {
		log.Fatalf("连接 Redis 失败: %v", err)
	}

	defer r.Close()

	e := echo.New()
	// 挂载 自定义的全局错误处理函数
	e.HTTPErrorHandler = response.HTTPErrorHandler
	// 挂载 恢复中间件
	e.Use(middleware.Recover())

	// 挂载 请求 ID 中间件
	e.Use(middleware.RequestIDWithConfig(middleware.RequestIDConfig{
		// 自定义请求 ID 生成器
		Generator: func() string {
			return uuid.Must(uuid.NewV7()).String()
		},
	}))

	// 使用 slog 记录日志：使用 JSON 格式化器，输出到标准输出
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	// 挂载 日志中间件
	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogStatus:    true, // 记录响应状态码
		LogURI:       true, // 记录请求 URI
		LogMethod:    true, // 记录请求方法
		HandleError:  true, // 把错误传递给 LogValuesFunc 处理
		LogRequestID: true, // 记录请求 ID
		LogLatency:   true, // 记录请求耗时
		// 自定义日志记录函数
		LogValuesFunc: func(c *echo.Context, v middleware.RequestLoggerValues) error {
			// 记录 正常 请求日志: 记录请求 URI、响应状态码、请求方法
			if v.Error == nil {
				logger.LogAttrs(c.Request().Context(), slog.LevelInfo, "REQUEST",
					slog.String("request_id", v.RequestID),
					slog.String("uri", v.URI),
					slog.String("method", v.Method),
					slog.Int("status", v.Status),
				)
			} else {
				// 记录 错误 日志: 记录请求 URI、响应状态码、错误信息
				logger.LogAttrs(c.Request().Context(), slog.LevelError, "REQUEST_ERROR", slog.String("uri", v.URI),
					slog.Int("status", v.Status), slog.String("error", v.Error.Error()), slog.String("request_id", v.RequestID))
			}
			return nil
		},
	}))

	// 挂载 路由
	health := &controller.Health{MySQL: sqlDB, Redis: r}                   // 健康检查控制器
	repo := repository.NewUserRepo(gdb)                                    // 用户仓库
	svc := service.NewUserService(repo, cfg.JWT.Secret, cfg.JWT.ExpiresIn) // 用户服务
	user := &controller.User{Svc: svc}                                     // 用户控制器

	// 创建 JWT 中间件
	jwtMW := authmw.JWT([]byte(cfg.JWT.Secret))

	// 注册路由集合: 健康检查、用户路由
	router.RegisterRouter(e, health, user, jwtMW)

	// 启动服务
	if err := e.Start(":" + cfg.Server.Port); err != nil {
		log.Fatal("start server: ", err)
	}
}
