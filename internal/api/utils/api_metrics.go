package utils

import "github.com/fajrikornel/go-todoapp/internal/metrics"

func IncrementApiSuccessMetric(api string) {
	metrics.IncrementCounter(
		ApiSuccess{
			Api: api,
		},
	)
}

func IncrementApiErrorMetric(api, errorCode string) {
	metrics.IncrementCounter(
		ApiError{
			Api:       api,
			ErrorCode: errorCode,
		},
	)
}

type ApiSuccess struct {
	metrics.Counter
	Api string
}

type ApiError struct {
	metrics.Counter
	Api       string
	ErrorCode string
}
