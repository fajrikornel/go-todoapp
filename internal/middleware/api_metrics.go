package middleware

import (
	"github.com/fajrikornel/go-todoapp/internal/metrics"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"regexp"
)

func MetricsMiddleware(handle httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		path := cleanPath(r.URL.Path)

		timer := metrics.NewTimer(
			ApiResponseTime{
				Path:   path,
				Method: r.Method,
			},
		)
		defer timer.ObserveDuration()

		metricResponseWriter := &MetricResponseWriter{
			ResponseWriter: w,
		}
		handle(metricResponseWriter, r, p)

		metrics.IncrementCounter(
			ApiHit{
				Path:         path,
				Method:       r.Method,
				ResponseCode: metricResponseWriter.statusCode,
			},
		)
	}
}

func cleanPath(path string) string {
	re := regexp.MustCompile("\\/\\d+")
	return re.ReplaceAllString(path, "/{...}")
}

type MetricResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (m *MetricResponseWriter) WriteHeader(statusCode int) {
	m.statusCode = statusCode
	m.ResponseWriter.WriteHeader(statusCode)
}

type ApiResponseTime struct {
	metrics.Histogram `help:"Histogram that tracks API response times from the server"`
	Path              string
	Method            string
}

type ApiHit struct {
	metrics.Counter `help:"Counter that tracks API hits to the server"`
	Path            string
	Method          string
	ResponseCode    int
}
