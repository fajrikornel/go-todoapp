package logging

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
)

type Log map[string]string

func Infof(ctx context.Context, format string, args ...interface{}) {
	log := make(Log)

	log.addLogLevel("INFO")
	if args == nil {
		log.addLogMessage(format)
	} else {
		log.addLogMessage(format, args...)
	}
	log.addCorrelationId(ctx)
	log.addTimestamp()

	log.printLog()
}

func Errorf(ctx context.Context, format string, args ...interface{}) {
	log := make(Log)

	log.addLogLevel("INFO")
	if args == nil {
		log.addLogMessage(format)
	} else {
		log.addLogMessage(format, args)
	}
	log.addCorrelationId(ctx)
	log.addTimestamp()

	log.printLog()
}

func (log *Log) addLogLevel(logLevel string) {
	(*log)["level"] = logLevel
}

func (log *Log) addLogMessage(format string, args ...interface{}) {
	if args != nil {
		(*log)["message"] = fmt.Sprintf(format, args...)
	} else {
		(*log)["message"] = format
	}

}

func (log *Log) addCorrelationId(ctx context.Context) {
	if correlationId := ctx.Value("correlationId"); correlationId != nil {
		(*log)["correlationId"] = correlationId.(string)
	}
}

func (log *Log) addTimestamp() {
	(*log)["timestamp"] = time.Now().Format(time.RFC3339)
}

func (log *Log) printLog() {
	logString, _ := json.Marshal(*log)
	fmt.Printf("%s\n", logString)
}
