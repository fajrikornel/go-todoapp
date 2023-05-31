package middleware

import (
	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
)

func LoggingMiddleware(handle httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		correlationId := uuid.NewString()
		log.Printf("Request received %s %s, correlationId: %s", r.Method, r.RequestURI, correlationId)
		handle(w, r, p)
		log.Printf("Request completed %s %s, correlationId: %s", r.Method, r.RequestURI, correlationId)
	}
}
