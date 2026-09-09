package server

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ChenHaoJie9527/Elk-Mall/internal/adaptor"
	"github.com/ChenHaoJie9527/Elk-Mall/internal/common/response"
	"github.com/ChenHaoJie9527/Elk-Mall/internal/config"
	"github.com/ChenHaoJie9527/Elk-Mall/internal/controller"
	"github.com/ChenHaoJie9527/Elk-Mall/internal/logger"
	authmw "github.com/ChenHaoJie9527/Elk-Mall/internal/middleware"
	"github.com/ChenHaoJie9527/Elk-Mall/internal/model/do"
	"github.com/ChenHaoJie9527/Elk-Mall/internal/repository"
	"github.com/ChenHaoJie9527/Elk-Mall/internal/router"
	"github.com/ChenHaoJie9527/Elk-Mall/internal/service"
	"github.com/google/uuid"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// Run 是进程组合根：加载配置、连基础设施、挂路由，然后阻塞到进程收到退出信号。
func Run() error {
	defer logger.Sync()

	cfg, err := loadConfig()
	if err != nil {
		return err
	}

	gdb, sqlDB, err := openMySQL(cfg)
	if err != nil {
		return err
	}
	// *sql.DB.Close 会关掉连接池；GORM 没有自己的 Close。
	defer sqlDB.Close()

	// 按 UserDO 建/补 users 表（开发阶段方便；生产一般改成独立迁移命令）
	if err := gdb.AutoMigrate(&do.UserDO{}); err != nil {
		return fmt.Errorf("自动迁移失败: %w", err)
	}

	rdb, err := openRedis(cfg)
	if err != nil {
		return err
	}
	defer rdb.Close()

	e := newEcho()
	registerApp(e, cfg, gdb, sqlDB, rdb)

	addr := ":" + cfg.Server.Port
	logger.Info("server starting", zap.String("app", cfg.App.Name), zap.String("addr", addr))
	return startServer(e, addr)
}

// loadConfig 有 ETCD_ADDR / -r 时读 etcd，否则读 -c 指定的本地 YAML（默认 config.yaml）。
func loadConfig() (*config.Config, error) {
	cfg, err := config.InitConfig()
	if err != nil {
		return nil, fmt.Errorf("加载配置失败: %w", err)
	}

	logger.Init(logger.Options{
		Env:        cfg.Server.Env,
		App:        cfg.App.Name,
		Level:      cfg.Server.LogLevel,
		File:       cfg.Server.LogFile,
		MaxSizeMB:  cfg.Server.LogMaxSizeMB,
		MaxBackups: cfg.Server.LogMaxBackups,
		MaxAgeDays: cfg.Server.LogMaxAgeDays,
	})

	logger.Info("config loaded",
		zap.String("source", config.Source()),
		zap.String("app", cfg.App.Name),
		zap.String("port", cfg.Server.Port),
		zap.String("env", cfg.Server.Env),
		zap.String("log_level", cfg.Server.LogLevel),
		zap.String("mysql_host", cfg.MySQL.Host),
		zap.String("redis_addr", cfg.Redis.Addr),
	)

	return cfg, nil
}

func openMySQL(cfg *config.Config) (*gorm.DB, *sql.DB, error) {
	gdb, err := adaptor.OpenMySql(cfg.MySQL)
	if err != nil {
		return nil, nil, fmt.Errorf("连接 MySQL 失败: %w", err)
	}

	sqlDB, err := gdb.DB()
	if err != nil {
		return nil, nil, fmt.Errorf("取出 MySQL 连接失败: %w", err)
	}
	return gdb, sqlDB, nil
}

func openRedis(cfg *config.Config) (*redis.Client, error) {
	rdb, err := adaptor.OpenRedis(cfg.Redis)
	if err != nil {
		return nil, fmt.Errorf("连接 Redis 失败: %w", err)
	}
	return rdb, nil
}

// newEcho 只负责创建引擎和挂「全局」中间件，不注册业务路由。
func newEcho() *echo.Echo {
	e := echo.NewWithConfig(echo.Config{
		Logger: logger.ToSlog(),
	})
	e.HTTPErrorHandler = response.HTTPErrorHandler // 业务 error → 统一 JSON
	e.Use(middleware.Recover())                    // panic 转 500，避免整个进程挂掉

	// 每个请求一个 ID，方便把多条日志串起来
	e.Use(middleware.RequestIDWithConfig(middleware.RequestIDConfig{
		Generator: func() string {
			return uuid.Must(uuid.NewV7()).String()
		},
	}))

	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogStatus:      true,
		LogURI:         true,
		LogMethod:      true,
		LogRequestID:   true,
		LogQueryParams: []string{"page", "lang"},
		LogLatency:     true, // 打开后 v.Latency 才有值，必须在下面真正打出来
		HandleError:    true, // 把 handler 返回的 error 传给 LogValuesFunc
		LogValuesFunc: func(c *echo.Context, v middleware.RequestLoggerValues) error {
			rec := logger.AccessRecord{
				Method:    v.Method,
				URI:       v.URI,
				Status:    v.Status,
				Latency:   v.Latency,
				RequestID: v.RequestID,
				Error:     v.Error,
			}
			logger.Access(rec)
			return nil
		},
	}))

	return e
}

// registerApp 组装依赖并注册路由：repo → service → controller → router。
func registerApp(e *echo.Echo, cfg *config.Config, gdb *gorm.DB, sqlDB *sql.DB, rdb *redis.Client) {
	health := &controller.Health{MySQL: sqlDB, Redis: rdb}
	repo := repository.NewUserRepo(gdb)
	svc := service.NewUserService(repo, cfg.JWT.Secret, cfg.JWT.ExpiresIn)
	user := &controller.User{Svc: svc}
	jwtMW := authmw.JWT([]byte(cfg.JWT.Secret))

	router.RegisterRouter(e, health, user, jwtMW)
}

// startServer 用 Echo v5 自带的优雅退出，不必再手写 v9 那种 signal + Shutdown 套娃。
//
// e.Start 内部其实就是：
//  1. signal.NotifyContext 监听 Ctrl+C / SIGTERM
//  2. ListenAndServe
//  3. 收到信号后 Shutdown，默认最多等 10 秒让进行中的请求结束
//
// 这里显式用 StartConfig，方便改超时、关掉默认 banner。
func startServer(e *echo.Echo, addr string) error {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	sc := echo.StartConfig{
		Address:         addr,
		HideBanner:      true,            // 隐藏 banner
		GracefulTimeout: 5 * time.Second, // 优雅退出超时时间 5s
	}
	return sc.Start(ctx, e)
}
