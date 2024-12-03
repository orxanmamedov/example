package logger

import (
	"context"

	"go.uber.org/zap"
)

type contextKey struct {
}

var (
	httpContextKey = &contextKey{}
)

func FromContext(ctx context.Context) *zap.SugaredLogger {
	l := global

	l = loggerWithHTTPContext(ctx, l)
	return l
}

func loggerWithHTTPContext(ctx context.Context, l *zap.SugaredLogger) *zap.SugaredLogger {
	if m, ok := ctx.Value(httpContextKey).(map[string]interface{}); ok {
		for k, v := range m {
			l = l.With(k, v)
		}
		return l

	}
	return l
}
