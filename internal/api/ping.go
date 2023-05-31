package api

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func PingHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.WriteHeader(200)
	w.Write([]byte("pong"))
}
