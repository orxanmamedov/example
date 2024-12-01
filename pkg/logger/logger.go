package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	global *zap.SugaredLogger
	level  = zap.NewAtomicLevelAt(zap.WarnLevel)
)

func init() {
	SetLogger(New(level))
}

func New(lvl zapcore.LevelEnabler, options ...zap.Option) *zap.SugaredLogger {
	if lvl == nil {
		lvl = level
	}
	sink := zapcore.AddSync(os.Stdout)
	options = append(options, zap.AddCallerSkip(1), zap.AddCaller())
	return zap.New(
		zapcore.NewCore(
			zapcore.NewJSONEncoder(zapcore.EncoderConfig{
				TimeKey:        "@timestamp",
				LevelKey:       "level",
				NameKey:        "logger",
				CallerKey:      "caller",
				MessageKey:     "message",
				StacktraceKey:  "stacktrace",
				LineEnding:     zapcore.DefaultLineEnding,
				EncodeLevel:    zapcore.LowercaseColorLevelEncoder,
				EncodeTime:     zapcore.ISO8601TimeEncoder,
				EncodeDuration: zapcore.SecondsDurationEncoder,
				EncodeCaller:   zapcore.ShortCallerEncoder,
			}),
			sink,
			lvl,
		),
		options...,
	).Sugar()
}

func Level() *zap.SugaredLogger {
	return global
}

func SetLogger(l *zap.SugaredLogger) {
	global = l
}
