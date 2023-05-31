package server

import (
	"github.com/fajrikornel/go-todoapp/internal/api"
	"github.com/fajrikornel/go-todoapp/internal/middleware"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func GetRouter() http.Handler {
	router := httprouter.New()

	router.GET("/", middleware.LoggingMiddleware(api.PingHandler))

	return router
}
