package logger

import (
	"fmt"
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"testing"
	"time"
)

func TestLog(t *testing.T) {
	Debug("debug")
	Info("info")
	Warn("warn")
	Error("error")
}

func TestSugarLog(t *testing.T) {
	Debugf("level: %s", "debug")
	Infof("level: %s", "info")
	Warnf("level: %s", "warn")
	Errorf("level: %s", "error")
}

func TestZap(t *testing.T) {
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder, // ISO8601 UTC 时间格式
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.FullCallerEncoder, // 全路径编码器
	}
	// 设置日志级别
	config := zap.Config{
		Level:            zap.NewAtomicLevelAt(zap.DebugLevel),                           // 日志级别
		Development:      true,                                                           // 开发模式，堆栈跟踪
		Encoding:         "json",                                                         // 输出格式 console 或 json
		EncoderConfig:    encoderConfig,                                                  // 编码器配置
		InitialFields:    map[string]interface{}{"serviceName": "merak-example_TestZap"}, // 初始化字段，如：添加一个服务器名称
		OutputPaths:      []string{"stdout", "merak-example.logger"},                     // 输出到指定文档 stdout（标准输出，正常颜色） stderr（错误输出，红色）
		ErrorOutputPaths: []string{"stderr"},
	}
	// 构建日志
	logger, err := config.Build()
	if err != nil {
		panic(fmt.Sprintf("logger 初始化失败: %v", err))
	}
	logger.Info("zap logger 初始化成功")
	logger.Info("无法获取网址",
		zap.String("url", "https://www.baidu.com"),
		zap.Int("attempt", 3),
		zap.Duration("backoff", time.Second),
	)
	logger.Error("错误", zap.String("error", "错误信息"))
}

func TestZapWithLumberjack(t *testing.T) {
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder, // ISO8601 UTC 时间格式
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.FullCallerEncoder, // 全路径编码器
	}
	hook := lumberjack.Logger{
		Filename:   "merak-example_TestZapWithLumberjack.logger", // 日志文档路径
		MaxSize:    128,                                          // 每个日志文档保存的最大尺寸 单位：M
		MaxBackups: 30,                                           // 日志文档最多保存多少个备份
		MaxAge:     7,                                            // 文档最多保存多少天
		Compress:   true,                                         // 是否压缩
	}
	// 设置日志级别
	atomicLevel := zap.NewAtomicLevel()
	atomicLevel.SetLevel(zap.DebugLevel)
	core := zapcore.NewCore(
		//zapcore.NewJSONEncoder(encoderConfig),                                           // 编码器配置
		zapcore.NewConsoleEncoder(encoderConfig),                                        // 编码器配置
		zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(&hook)), // 打印到控制面板和文档
		atomicLevel, // 日志级别
	)
	caller := zap.AddCaller()                                                               // 开启开发模式，堆栈跟踪
	development := zap.Development()                                                        // 开启文档及行号
	fields := zap.Fields(zap.String("serviceName", "merak-example"))                        // 设置初始化字段
	logger := zap.New(core, caller, development, fields, zap.AddStacktrace(zap.ErrorLevel)) // 构造日志

	logger.Info("logger 初始化成功")
	logger.Info("无法获取网址",
		zap.String("url", "https://www.baidu.com"),
		zap.Int("attempt", 3),
		zap.Duration("backoff", time.Second))
	logger.Error("错误", zap.String("error", "错误信息"))
}
