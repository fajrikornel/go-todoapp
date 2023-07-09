package logging

import (
	"context"
	"go.uber.org/zap"
)

type ZapLogger struct {
	sl *zap.SugaredLogger
}

func NewZapLogger() *ZapLogger {
	z, _ := zap.NewProduction()
	s := z.Sugar()

	return &ZapLogger{
		sl: s,
	}
}

func (zl *ZapLogger) infof(ctx context.Context, format string, args ...interface{}) {
	correlationId := extractCorrelationIdFromContext(ctx)
	if args == nil {
		zl.sl.With("correlationId", correlationId).Infof(format)
	} else {
		zl.sl.With("correlationId", correlationId).Infof(format, args...)
	}
}

func (zl *ZapLogger) errorf(ctx context.Context, format string, args ...interface{}) {
	correlationId := extractCorrelationIdFromContext(ctx)
	if args == nil {
		zl.sl.With("correlationId", correlationId).Errorf(format)
	} else {
		zl.sl.With("correlationId", correlationId).Errorf(format, args...)
	}
}

func extractCorrelationIdFromContext(ctx context.Context) *string {
	if correlationId := ctx.Value("correlationId"); correlationId != nil {
		correlationIdString := correlationId.(string)
		return &correlationIdString
	}
	return nil
}
