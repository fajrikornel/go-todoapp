package logging

import (
	"context"
	"github.com/fajrikornel/go-todoapp/internal/config"
)

type Logger interface {
	infof(ctx context.Context, format string, args ...interface{})
	errorf(ctx context.Context, format string, args ...interface{})
}

var logger Logger

func init() {
	zapLoggerEnabled := config.GetEnableZapLogger()
	if zapLoggerEnabled {
		logger = NewZapLogger()
	} else {
		logger = NewDefaultLogger()
	}
}

func Infof(ctx context.Context, format string, args ...interface{}) {
	if args == nil {
		logger.infof(ctx, format)
	} else {
		logger.infof(ctx, format, args...)
	}
}

func Errorf(ctx context.Context, format string, args ...interface{}) {
	if args == nil {
		logger.errorf(ctx, format)
	} else {
		logger.errorf(ctx, format, args...)
	}
}
