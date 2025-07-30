package pkg

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"path/filepath"
)

// Logger 全局日志实例
var Logger *zap.Logger

// InitLogger 初始化zap日志，配置不同级别日志的存储位置
func InitLogger() error {
	// 创建日志目录
	errLogPath := filepath.Join("internal/basic/zap/zaperr/err.log")
	runLogPath := filepath.Join("internal/basic/zap/zaprun/run.log")

	// 确保目录存在
	if err := os.MkdirAll(filepath.Dir(errLogPath), 0755); err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(runLogPath), 0755); err != nil {
		return err
	}

	// 创建错误日志文件写入器
	errFile, err := os.OpenFile(errLogPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	// 创建运行日志文件写入器
	runFile, err := os.OpenFile(runLogPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	// 配置编码器
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "timestamp",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05"),
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	// 创建编码器
	encoder := zapcore.NewJSONEncoder(encoderConfig)

	// 错误级别过滤器 (ERROR, FATAL, PANIC)
	errorLevelEnabler := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.ErrorLevel
	})

	// 信息级别过滤器 (INFO, DEBUG)
	infoLevelEnabler := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.InfoLevel && lvl < zapcore.ErrorLevel
	})

	// 创建核心
	core := zapcore.NewTee(
		zapcore.NewCore(encoder, zapcore.AddSync(errFile), errorLevelEnabler),    // 错误日志写入errzap/log
		zapcore.NewCore(encoder, zapcore.AddSync(runFile), infoLevelEnabler),     // 运行日志写入runzap/log
		zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), zapcore.DebugLevel), // 控制台输出
	)

	// 创建logger
	Logger = zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
	return nil
}

// GetLogger 获取日志实例
func GetLogger() *zap.Logger {
	if Logger == nil {
		// 如果Logger未初始化，使用默认配置
		if err := InitLogger(); err != nil {
			Logger, _ = zap.NewProduction()
		}
	}
	return Logger
}

// LogError 记录错误日志到errzap/log
func LogError(msg string, fields ...zap.Field) {
	GetLogger().Error(msg, fields...)
}

// LogInfo 记录信息日志到runzap/log
func LogInfo(msg string, fields ...zap.Field) {
	GetLogger().Info(msg, fields...)
}

// Sync 同步日志缓冲区
func Sync() {
	if Logger != nil {
		Logger.Sync()
	}
}
