package logger

import (
	"context"
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

func SetLevel(l string) {
	var zapLevel zapcore.Level
	if err := zapLevel.Set(l); err != nil {
		Warnf(context.Background(), "failed  parse log level %#q: %s, setting to WARN", level, err)
		zapLevel = zapcore.WarnLevel
	}
}

func Logger() *zap.SugaredLogger {
	return global
}

func SetLogger(l *zap.SugaredLogger) {
	global = l
}

func Debug(ctx context.Context, args ...interface{}) {
	FromContext(ctx).Debug(args...)
}

func Debugf(ctx context.Context, format string, args ...interface{}) {
	FromContext(ctx).Debugf(format, args...)
}

func DebugKV(ctx context.Context, message string, kvs ...interface{}) {
	FromContext(ctx).Debugw(message, kvs...)
}

func Info(ctx context.Context, args ...interface{}) {
	FromContext(ctx).Info(args...)
}

func Infof(ctx context.Context, format string, args ...interface{}) {
	FromContext(ctx).Infof(format, args...)
}

func InfoKV(ctx context.Context, message string, kvs ...interface{}) {
	FromContext(ctx).Infow(message, kvs...)
}

func Warn(ctx context.Context, args ...interface{}) {
	FromContext(ctx).Warn(args...)
}

func Warnf(ctx context.Context, format string, args ...interface{}) {
	FromContext(ctx).Warnf(format, args...)
}

func WarnKV(ctx context.Context, message string, kvs ...interface{}) {
	FromContext(ctx).Warnw(message, kvs...)
}

func Error(ctx context.Context, args ...interface{}) {
	FromContext(ctx).Error(args...)
}

func Errorf(ctx context.Context, format string, args ...interface{}) {
	FromContext(ctx).Errorf(format, args...)
}

func ErrorKV(ctx context.Context, message string, kvs ...interface{}) {
	FromContext(ctx).Errorw(message, kvs...)
}

func Fatal(ctx context.Context, args ...interface{}) {
	FromContext(ctx).Fatal(args...)
}

func Fatalf(ctx context.Context, format string, args ...interface{}) {
	FromContext(ctx).Fatalf(format, args...)
}

func FatalKV(ctx context.Context, message string, kvs ...interface{}) {
	FromContext(ctx).Fatalw(message, kvs...)
}
