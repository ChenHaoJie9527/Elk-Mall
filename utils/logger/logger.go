package logger

import (
	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	logger *zap.Logger                            // zap 日志记录器
	atom   = zap.NewAtomicLevelAt(zap.DebugLevel) // 日志级别
)

func init() {
	config := zap.Config{
		Level:       atom,   // 日志级别
		Development: false,  // 是否开发模式
		Encoding:    "json", // 指定 JSON 编码
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey: "msg",    // 消息键
			LevelKey:   "level",  // 级别键
			TimeKey:    "time",   // 时间键
			CallerKey:  "caller", // 调用者键
			//StacktraceKey: "stacktrace",
			EncodeTime:   zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05.000"), // 时间编码
			EncodeLevel:  zapcore.LowercaseLevelEncoder,                          // 级别编码
			EncodeCaller: zapcore.ShortCallerEncoder,                             // 调用者编码
		},
		OutputPaths:      []string{"stdout"}, // 输出路径
		ErrorOutputPaths: []string{"stderr"}, // 错误输出路径
	}

	// 构建日志记录器
	tempLogger, err := config.Build()
	if err != nil {
		panic(err)
	}

	// 将日志记录器赋值给全局变量
	// 添加调用者、调用者跳过1级、添加堆栈跟踪
	logger = tempLogger.WithOptions(zap.AddCaller(), zap.AddCallerSkip(1), zap.AddStacktrace(zap.ErrorLevel))

}

func SetLevel(level string) {
	// 解析日志级别
	tLevel, err := zapcore.ParseLevel(level)
	if err != nil {
		fmt.Printf("invalid level, input: %s", level)
		return
	}
	atom.SetLevel(tLevel)
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
