package server

import (
	"github.com/fajrikornel/go-todoapp/internal/api"
	"github.com/fajrikornel/go-todoapp/internal/api/v1"
	"github.com/fajrikornel/go-todoapp/internal/middleware"
	"github.com/fajrikornel/go-todoapp/internal/repository"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func GetRouter(repository repository.ProjectRepository) http.Handler {
	router := httprouter.New()

	router.GET("/", middleware.LoggingMiddleware(api.PingHandler))
	router.POST("/v1/projects", middleware.LoggingMiddleware(v1.CreateProjectHandler(repository)))

	return router
}
