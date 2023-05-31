package server

import (
	"github.com/fajrikornel/go-todoapp/internal/api"
	v1 "github.com/fajrikornel/go-todoapp/internal/api/v1"
	"github.com/fajrikornel/go-todoapp/internal/db"
	"github.com/fajrikornel/go-todoapp/internal/middleware"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func GetRouter(store *db.SqlStore) http.Handler {
	router := httprouter.New()

	router.GET("/", middleware.LoggingMiddleware(api.PingHandler))
	router.POST("/v1/projects", middleware.LoggingMiddleware(v1.CreateProjectHandler(store)))

	return router
}
