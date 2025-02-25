package logger

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Config 日志配置
type Config struct {
	// 日志文件路径
	Filename string
	// 单个日志文件最大大小（MB）
	MaxSize int
	// 保留的旧日志文件个数
	MaxBackups int
	// 保留的日志天数
	MaxAge int
	// 是否压缩旧日志
	Compress bool
	// 日志级别
	Level string
	// 是否输出到控制台
	Console bool
}

// NewRotateLogger 创建带轮转功能的日志记录器
func NewRotateLogger(config Config) (*zap.Logger, error) {
	// 确保日志目录存在
	if err := os.MkdirAll(filepath.Dir(config.Filename), 0755); err != nil {
		return nil, fmt.Errorf("create log directory failed: %w", err)
	}

	// 设置日志级别
	level := zap.InfoLevel
	if err := level.UnmarshalText([]byte(config.Level)); err != nil {
		return nil, fmt.Errorf("parse log level failed: %w", err)
	}

	// 创建轮转日志写入器
	writer := &lumberjack.Logger{
		Filename:   config.Filename,
		MaxSize:    config.MaxSize,    // MB
		MaxBackups: config.MaxBackups, // 文件个数
		MaxAge:     config.MaxAge,     // 天数
		Compress:   config.Compress,   // 是否压缩
	}

	// 配置编码器
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	// 创建核心
	var core zapcore.Core
	if config.Console {
		// 同时输出到文件和控制台
		consoleEncoder := zapcore.NewConsoleEncoder(encoderConfig)
		fileEncoder := zapcore.NewJSONEncoder(encoderConfig)
		core = zapcore.NewTee(
			zapcore.NewCore(fileEncoder, zapcore.AddSync(writer), level),
			zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), level),
		)
	} else {
		// 只输出到文件
		fileEncoder := zapcore.NewJSONEncoder(encoderConfig)
		core = zapcore.NewCore(fileEncoder, zapcore.AddSync(writer), level)
	}

	// 创建日志记录器
	logger := zap.New(core,
		zap.AddCaller(),
		zap.AddCallerSkip(1),
		zap.AddStacktrace(zapcore.ErrorLevel),
	)

	return logger, nil
}

// GetDefaultConfig 获取默认配置
func GetDefaultConfig() Config {
	return Config{
		Filename:   "logs/app.log",
		MaxSize:    100,    // 100MB
		MaxBackups: 30,     // 保留30个备份
		MaxAge:     7,      // 保留7天
		Compress:   true,   // 压缩旧日志
		Level:      "info", // 默认info级别
		Console:    true,   // 默认同时输出到控制台
	}
}

// CleanOldLogs 清理过期日志
func CleanOldLogs(dir string, maxAge time.Duration) error {
	now := time.Now()
	return filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && info.ModTime().Add(maxAge).Before(now) {
			if err := os.Remove(path); err != nil {
				return fmt.Errorf("remove old log file failed: %w", err)
			}
		}
		return nil
	})
}
