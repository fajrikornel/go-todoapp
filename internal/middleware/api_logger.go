package middleware

import (
	"context"
	"github.com/fajrikornel/go-todoapp/internal/logging"
	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func LoggingMiddleware(handle httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		correlationId := uuid.NewString()
		correlationIdCtx := context.WithValue(r.Context(), "correlationId", correlationId)
		reqWithCorrelation := r.WithContext(correlationIdCtx)

		logging.Infof(correlationIdCtx, "Request received %s %s", r.Method, r.RequestURI)
		handle(w, reqWithCorrelation, p)
		logging.Infof(correlationIdCtx, "Request completed %s %s", r.Method, r.RequestURI)
	}
}
