package server

import (
	"github.com/fajrikornel/go-todoapp/internal/api"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func GetRouter() http.Handler {
	router := httprouter.New()

	router.GET("/", api.PingHandler)

	return router
}
