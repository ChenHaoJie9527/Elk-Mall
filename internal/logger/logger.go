package logger

import (
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"strings"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/exp/zapslog"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	// 包级 zap 实例；业务侧用 Debug/Info/Warn/Error，不要直接拿这个变量
	logger *zap.Logger
	// 原子级别：启动后可用 SetLevel 改，不必重建 logger
	atom = zap.NewAtomicLevelAt(zap.DebugLevel)
	// local 用彩色 key=value；prod 仍走 JSON
	prettyConsole bool
)

type Options struct {
	Env        string
	App        string
	Level      string
	File       string
	MaxSizeMB  int
	MaxBackups int
	MaxAgeDays int
}

func Init(opts Options) {
	env := strings.ToLower(strings.TrimSpace(opts.Env))
	if env == "" {
		env = "local"
	}

	var built *zap.Logger
	var err error

	if env == "local" {
		// 开发环境
		built, err = newDevLogger()
	} else {
		built, err = newProdLogger(opts)
	}

	if err != nil {
		panic(err)
	}

	prettyConsole = env == "local"

	zopts := []zap.Option{
		zap.AddCaller(),
		zap.AddCallerSkip(1),
		zap.AddStacktrace(zap.ErrorLevel),
	}
	// 本地每行都带 app/env 很吵；生产 JSON 留给采集用
	if !prettyConsole {
		zopts = append(zopts,
			zap.Fields(zap.String("app", opts.App), zap.String("env", env)),
		)
	}
	logger = built.WithOptions(zopts...)

	SetLevel(opts.Level)
}

func newDevLogger() (*zap.Logger, error) {
	encCfg := zapcore.EncoderConfig{
		TimeKey:          "time",
		LevelKey:         "level",
		CallerKey:        "caller",
		MessageKey:       "msg",
		StacktraceKey:    "stacktrace",
		LineEnding:       zapcore.DefaultLineEnding,
		EncodeLevel:      zapcore.CapitalColorLevelEncoder,
		EncodeTime:       zapcore.TimeEncoderOfLayout("15:04:05.000"),
		EncodeDuration:   zapcore.StringDurationEncoder,
		EncodeCaller:     zapcore.ShortCallerEncoder,
		ConsoleSeparator: "  ",
	}
	core := zapcore.NewCore(newKVEncoder(encCfg), zapcore.AddSync(os.Stdout), atom)
	return zap.New(core), nil
}

func newProdLogger(opts Options) (*zap.Logger, error) {
	file := opts.File
	if file == "" {
		file = "./logs/server.log"
	}
	rolling := &lumberjack.Logger{
		Filename:   file,
		MaxSize:    maxInt(opts.MaxSizeMB, 100),
		MaxBackups: maxInt(opts.MaxBackups, 3),
		MaxAge:     maxInt(opts.MaxAgeDays, 28),
		Compress:   true,
	}
	encCfg := zap.NewProductionEncoderConfig()
	encCfg.TimeKey = "time"
	encCfg.EncodeTime = zapcore.ISO8601TimeEncoder
	encCfg.EncodeLevel = zapcore.LowercaseLevelEncoder
	encCfg.EncodeCaller = zapcore.ShortCallerEncoder
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encCfg),
		zapcore.NewMultiWriteSyncer(
			zapcore.AddSync(rolling),
			zapcore.AddSync(os.Stdout),
		),
		atom,
	)
	return zap.New(core), nil
}

func maxInt(v, fallback int) int {
	if v <= 0 {
		return fallback
	}
	return v
}

// SetLevel 按配置热改日志级别，例如 debug / info / warn / error。
func SetLevel(level string) {
	level = strings.TrimSpace(level)
	if level == "" {
		return
	}
	lev, err := zapcore.ParseLevel(level)
	if err != nil {
		fmt.Printf("invalid log level: %s\n", level)
		return
	}
	atom.SetLevel(lev)
}

func Debug(msg string, fields ...zap.Field) {
	logger.Debug(msg, fields...)
}

func Info(msg string, fields ...zap.Field) {
	logger.Info(msg, fields...)
}

func Warn(msg string, fields ...zap.Field) {
	logger.Warn(msg, fields...)
}

func Error(msg string, fields ...zap.Field) {
	logger.Error(msg, fields...)
}

// AccessRecord 一条 HTTP 访问日志。本地打成单行，生产仍是结构化字段。
type AccessRecord struct {
	Method    string
	URI       string
	Status    int
	Latency   time.Duration
	RequestID string
	Error     error
}

func Access(r AccessRecord) {
	fields := []zap.Field{
		zap.String("method", r.Method),
		zap.String("uri", r.URI),
		zap.Int("status", r.Status),
		zap.Duration("latency", r.Latency),
	}
	if r.RequestID != "" {
		fields = append(fields, zap.String("request_id", r.RequestID))
	}
	l := logger.WithOptions(zap.WithCaller(false))
	if r.Error != nil {
		l.Error("request", append(fields, zap.Error(r.Error))...)
		return
	}
	l.Info("request", fields...)
}

func formatLatency(d time.Duration) string {
	switch {
	case d < time.Millisecond:
		return strconv.FormatInt(d.Microseconds(), 10) + "µs"
	case d < time.Second:
		ms := float64(d) / float64(time.Millisecond)
		switch {
		case ms >= 100:
			return strconv.FormatFloat(ms, 'f', 0, 64) + "ms"
		case ms >= 10:
			return strconv.FormatFloat(ms, 'f', 1, 64) + "ms"
		default:
			return strconv.FormatFloat(ms, 'f', 2, 64) + "ms"
		}
	default:
		return strconv.FormatFloat(d.Seconds(), 'f', 2, 64) + "s"
	}
}

func colorStatus(code int) string {
	s := strconv.Itoa(code)
	switch {
	case code >= 500:
		return "\x1b[31m" + s + "\x1b[0m"
	case code >= 400:
		return "\x1b[33m" + s + "\x1b[0m"
	default:
		return "\x1b[32m" + s + "\x1b[0m"
	}
}

// Sync 把缓冲里的日志刷出去。进程退出前调用；stdout 上失败可以忽略。
func Sync() {
	_ = logger.Sync()
}

// Slog 把 zap 转成 Echo v5 需要的 *slog.Logger，启动/关闭日志也走同一套 JSON。
func ToSlog() *slog.Logger {
	// 将 zap 的 Core 转换为 slog 的 Handler
	h := zapslog.NewHandler(logger.Core(), zapslog.WithCaller(true))
	// 创建 slog 的 Logger
	l := slog.New(h)
	return l
}
