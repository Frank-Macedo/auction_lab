package logger

import (
	"github.com/pingcap/log"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	Log *zap.Logger
)

func init() {
	logConfiguration := zap.Config{
		Level:    zap.NewAtomicLevelAt(zap.InfoLevel),
		Encoding: "json",
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey:   "message",
			LevelKey:     "level",
			TimeKey:      "time",
			EncodeTime:   zapcore.ISO8601TimeEncoder,
			EncodeLevel:  zapcore.LowercaseColorLevelEncoder,
			EncodeCaller: zapcore.ShortCallerEncoder,
		},
	}

	Log, _ = logConfiguration.Build()
}

func Info(message string, tags ...zap.Field) {
	Log.Info(message, tags...)
	log.Sync()
}

func Error(message string, err error, tags ...zap.Field) {
	tags = append(tags, zap.NamedError("error", err))
	Log.Error(message, tags...)
	log.Sync()
}

func Debug(message string, fields ...zap.Field) {

	Log.Debug(message, fields...)
	log.Sync()
}

func Warn(message string, fields ...zap.Field) {
	Log.Warn(message, fields...)
	log.Sync()
}
