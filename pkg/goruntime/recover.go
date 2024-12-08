package goruntime

import (
	"context"
	"example/pkg/logger"
	"fmt"
	"runtime/debug"
	"strings"
)

func RecoverPanic(ctx context.Context, name string) {
	if p := recover(); p != nil {
		
	}
}

func HandlePanic(ctx context.Context, name string, p interface{}) {
	if p != nil {
		stack := strings.Split(string(debug.Stack()),"\n")
		logger.ErrorKV(ctx, fmt.Sprintf("%s %v", name, p), "stack_trace", stack)
	}
}
