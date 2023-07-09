package logging

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
)

type DefaultLog map[string]string

type DefaultLogger struct{}

func NewDefaultLogger() *DefaultLogger {
	return &DefaultLogger{}
}

func (dl *DefaultLogger) infof(ctx context.Context, format string, args ...interface{}) {
	log := make(DefaultLog)

	log.defaultAddLogLevel("INFO")
	if args == nil {
		log.defaultAddLogMessage(format)
	} else {
		log.defaultAddLogMessage(format, args...)
	}
	log.defaultAddCorrelationId(ctx)
	log.defaultAddTimestamp()

	log.defaultPrintLog()
}

func (dl *DefaultLogger) errorf(ctx context.Context, format string, args ...interface{}) {
	log := make(DefaultLog)

	log.defaultAddLogLevel("INFO")
	if args == nil {
		log.defaultAddLogMessage(format)
	} else {
		log.defaultAddLogMessage(format, args)
	}
	log.defaultAddCorrelationId(ctx)
	log.defaultAddTimestamp()

	log.defaultPrintLog()
}

func (log *DefaultLog) defaultAddLogLevel(logLevel string) {
	(*log)["level"] = logLevel
}

func (log *DefaultLog) defaultAddLogMessage(format string, args ...interface{}) {
	if args != nil {
		(*log)["message"] = fmt.Sprintf(format, args...)
	} else {
		(*log)["message"] = format
	}

}

func (log *DefaultLog) defaultAddCorrelationId(ctx context.Context) {
	if correlationId := ctx.Value("correlationId"); correlationId != nil {
		(*log)["correlationId"] = correlationId.(string)
	}
}

func (log *DefaultLog) defaultAddTimestamp() {
	(*log)["timestamp"] = time.Now().Format(time.RFC3339)
}

func (log *DefaultLog) defaultPrintLog() {
	logString, _ := json.Marshal(*log)
	fmt.Printf("%s\n", logString)
}
