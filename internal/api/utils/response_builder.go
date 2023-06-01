package utils

import (
	"context"
	"encoding/json"
	"github.com/fajrikornel/go-todoapp/internal/logging"
	"net/http"
)

func ReturnErrorResponse(ctx context.Context, w http.ResponseWriter, httpCode int, responseBody interface{}, err error) {
	logging.Errorf(ctx, "Returned error response: %v", err.Error())
	w.WriteHeader(httpCode)
	json.NewEncoder(w).Encode(responseBody)
}

func ReturnSuccessResponse(ctx context.Context, w http.ResponseWriter, responseBody interface{}) {
	logging.Infof(ctx, "Returned success response: %v", responseBody)
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(responseBody)
}
