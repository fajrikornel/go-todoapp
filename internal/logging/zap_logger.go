package logging

import (
	"context"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"log"
	"os"
	"time"
)

type ZapLogger struct {
	sl *zap.SugaredLogger
}

func NewZapLogger() *ZapLogger {
	pe := zap.NewProductionEncoderConfig()
	pe.EncodeTime = zapcore.ISO8601TimeEncoder

	fileEncoder := zapcore.NewJSONEncoder(pe)
	consoleEncoder := zapcore.NewConsoleEncoder(pe)

	_ = os.Mkdir("logs/", 0751)
	logName := "logs/" + time.Now().UTC().Format(time.RFC3339) + ".log"
	file, err := os.OpenFile(logName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0751)
	if err != nil {
		log.Fatalf("FAILED INITIALIZING LOGGER: %s", err.Error())
	}

	core := zapcore.NewTee(
		zapcore.NewCore(fileEncoder, zapcore.AddSync(file), zap.DebugLevel),
		zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), zap.DebugLevel),
	)

	z := zap.New(core)
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
