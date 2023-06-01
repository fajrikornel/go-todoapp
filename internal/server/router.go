package server

import (
	"github.com/fajrikornel/go-todoapp/internal/api"
	"github.com/fajrikornel/go-todoapp/internal/api/v1"
	"github.com/fajrikornel/go-todoapp/internal/middleware"
	"github.com/fajrikornel/go-todoapp/internal/repository"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func GetRouter(pRepository repository.ProjectRepository, iRepository repository.ItemRepository) http.Handler {
	router := httprouter.New()

	router.GET("/", middleware.LoggingMiddleware(api.PingHandler))
	router.POST("/v1/projects", middleware.LoggingMiddleware(v1.CreateProjectHandler(pRepository)))
	router.POST("/v1/projects/:projectId", middleware.LoggingMiddleware(v1.CreateItemHandler(iRepository)))
	router.GET("/v1/projects/:projectId", middleware.LoggingMiddleware(v1.GetProjectHandler(pRepository)))
	router.GET("/v1/projects/:projectId/:itemId", middleware.LoggingMiddleware(v1.GetItemHandler(iRepository)))
	router.PATCH("/v1/projects/:projectId", middleware.LoggingMiddleware(v1.UpdateProjectHandler(pRepository)))
	router.PATCH("/v1/projects/:projectId/:itemId", middleware.LoggingMiddleware(v1.UpdateItemHandler(iRepository)))

	return router
}
