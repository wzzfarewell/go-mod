package logger

import (
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"strings"
)

var (
	Logger        *zap.Logger
	DefaultConfig = Config{
		Level:       "debug",
		Encoding:    "json",
		Development: false,
		FileName:    "logs/infrastructure.log",
		MaxSize:     128,
		MaxBackups:  30,
		MaxAge:      7,
		Compress:    false,
	}
)

type Config struct {
	Level       string // logger level: debug, info, warn, error, panic, fatal
	Encoding    string // logger encoding: json, console
	Development bool   // development mode
	FileName    string // logger file name
	MaxSize     int    // logger file max size, unit: MB
	MaxBackups  int    // logger file max backups
	MaxAge      int    // logger file max age, unit: day
	Compress    bool   // Compress determines if the rotated logger files should be compressed using gzip.
}

func init() {
	Init(DefaultConfig)
}

func Init(cfg Config) {
	Logger = newLogger(cfg)
}

func Debug(msg string, fields ...zap.Field) {
	Logger.Debug(msg, fields...)
}

func Info(msg string, fields ...zap.Field) {
	Logger.Info(msg, fields...)
}

func Warn(msg string, fields ...zap.Field) {
	Logger.Warn(msg, fields...)
}

func Error(msg string, fields ...zap.Field) {
	Logger.Error(msg, fields...)
}

func Panic(msg string, fields ...zap.Field) {
	Logger.Panic(msg, fields...)
}

func Fatal(msg string, fields ...zap.Field) {
	Logger.Fatal(msg, fields...)
}

func Debugf(format string, args ...any) {
	Logger.Sugar().Debugf(format, args...)
}

func Infof(format string, args ...any) {
	Logger.Sugar().Infof(format, args...)
}

func Warnf(format string, args ...any) {
	Logger.Sugar().Warnf(format, args...)
}

func Errorf(format string, args ...any) {
	Logger.Sugar().Errorf(format, args...)
}

func Panicf(format string, args ...any) {
	Logger.Sugar().Panicf(format, args...)
}

func Fatalf(format string, args ...any) {
	Logger.Sugar().Fatalf(format, args...)
}

func newLogger(cfg Config) *zap.Logger {
	hook := lumberjack.Logger{
		Filename:   cfg.FileName,
		MaxSize:    cfg.MaxSize,
		MaxBackups: cfg.MaxBackups,
		MaxAge:     cfg.MaxAge,
		Compress:   cfg.Compress,
	}
	atomicLevel := zap.NewAtomicLevel()
	atomicLevel.SetLevel(zapLogLevel(cfg.Level))
	core := zapcore.NewCore(
		zapEncoder(cfg.Encoding),
		zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(&hook)),
		atomicLevel,
	)
	caller := zap.AddCaller()
	development := zap.Development()
	fields := zap.Fields()
	return zap.New(core, caller, development, fields, zap.AddStacktrace(zap.ErrorLevel), zap.AddCallerSkip(1))
}

func zapEncoder(encoding string) zapcore.Encoder {
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
	if strings.EqualFold(encoding, "console") {
		return zapcore.NewConsoleEncoder(encoderConfig)
	}
	return zapcore.NewJSONEncoder(encoderConfig)
}

func zapLogLevel(level string) zapcore.Level {
	switch level {
	case "debug", "DEBUG":
		return zap.DebugLevel
	case "info", "INFO":
		return zap.InfoLevel
	case "warn", "WARN":
		return zap.WarnLevel
	case "error", "ERROR":
		return zap.ErrorLevel
	case "panic", "PANIC":
		return zap.PanicLevel
	case "fatal", "FATAL":
		return zap.FatalLevel
	default:
		return zap.InfoLevel
	}
}
