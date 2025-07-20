package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Log *zap.Logger

func Init(logLevel string) {
	var zapLevel zapcore.Level
	if err := zapLevel.UnmarshalText([]byte(logLevel)); err != nil {
		zapLevel = zapcore.InfoLevel // 默认级别
	}

	config := zap.Config{
		Level:            zap.NewAtomicLevelAt(zapLevel),
		Development:      false,
		Encoding:         "json",
		EncoderConfig:    zap.NewProductionEncoderConfig(),
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
	}

	var err error
	Log, err = config.Build()
	if err != nil {
		panic(err)
	}

	zap.ReplaceGlobals(Log)
}

func Sync() {
	_ = Log.Sync()
}

// 添加全局快捷方法
func Fatal(msg string, fields ...zap.Field) {
	Log.Fatal(msg, fields...)
}
func Info(msg string, fields ...zap.Field) {
	Log.Info(msg, fields...)
}

func Error(msg string, fields ...zap.Field) {
	Log.Error(msg, fields...)
}
func Debug(msg string, fields ...zap.Field) {
	Log.Debug(msg, fields...)
}
func Warn(msg string, fields ...zap.Field) {
	Log.Warn(msg, fields...)
}
