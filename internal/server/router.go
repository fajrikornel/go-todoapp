package server

import (
	"github.com/fajrikornel/go-todoapp/internal/api"
	"github.com/fajrikornel/go-todoapp/internal/api/v1"
	"github.com/fajrikornel/go-todoapp/internal/metrics"
	m "github.com/fajrikornel/go-todoapp/internal/middleware"
	"github.com/fajrikornel/go-todoapp/internal/repository"
	"github.com/julienschmidt/httprouter"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)

func GetRouter(pRepository repository.ProjectRepository, iRepository repository.ItemRepository) http.Handler {
	router := httprouter.New()

	router.GET("/", m.MetricsMiddleware(m.LoggingMiddleware(api.PingHandler)))
	router.POST("/v1/projects", m.MetricsMiddleware(m.LoggingMiddleware(v1.CreateProjectHandler(pRepository))))
	router.POST("/v1/projects/:projectId", m.MetricsMiddleware(m.LoggingMiddleware(v1.CreateItemHandler(iRepository))))
	router.GET("/v1/projects/:projectId", m.MetricsMiddleware(m.LoggingMiddleware(v1.GetProjectHandler(pRepository))))
	router.GET("/v1/projects/:projectId/:itemId", m.MetricsMiddleware(m.LoggingMiddleware(v1.GetItemHandler(iRepository))))
	router.PATCH("/v1/projects/:projectId", m.MetricsMiddleware(m.LoggingMiddleware(v1.UpdateProjectHandler(pRepository))))
	router.PATCH("/v1/projects/:projectId/:itemId", m.MetricsMiddleware(m.LoggingMiddleware(v1.UpdateItemHandler(iRepository))))
	router.DELETE("/v1/projects/:projectId", m.MetricsMiddleware(m.LoggingMiddleware(v1.DeleteProjectHandler(pRepository))))
	router.DELETE("/v1/projects/:projectId/:itemId", m.MetricsMiddleware(m.LoggingMiddleware(v1.DeleteItemHandler(iRepository))))

	router.Handler("GET", "/metrics", promhttp.HandlerFor(
		metrics.GetGatherer(),
		promhttp.HandlerOpts{
			Registry: metrics.GetRegisterer(),
		}),
	)

	return router
}
