package utils

import (
	"encoding/json"
	"log"
	"net/http"
)

func ReturnErrorResponse(w http.ResponseWriter, httpCode int, responseBody interface{}, err error) {
	log.Printf("Returned error response: %v", err.Error())
	w.WriteHeader(httpCode)
	json.NewEncoder(w).Encode(responseBody)
}

func ReturnSuccessResponse(w http.ResponseWriter, responseBody interface{}) {
	log.Printf("Returned success response: %+v", responseBody)
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(responseBody)
}
