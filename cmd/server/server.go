package main

import (
	"context"
	"fmt"
	"github.com/fajrikornel/go-todoapp/internal/config"
	"github.com/fajrikornel/go-todoapp/internal/db"
	"github.com/fajrikornel/go-todoapp/internal/logging"
	"github.com/fajrikornel/go-todoapp/internal/repository"
	"github.com/fajrikornel/go-todoapp/internal/server"
	"net/http"
)

func main() {
	dbConfig := config.GetDbConfig()
	store, err := db.GetSqlStore(&dbConfig)
	if err != nil {
		logging.Errorf(context.Background(), "ERROR INITIALIZING DB: {}", err.Error())
		return
	}

	projectRepository := repository.NewProjectRepository(store)
	itemRepository := repository.NewItemRepository(store)
	router := server.GetRouter(projectRepository, itemRepository)

	addr := fmt.Sprintf(":%d", config.GetAppPort())

	logging.Infof(context.Background(), "STARTING WEB SERVER")
	logging.Errorf(context.Background(), "%s", http.ListenAndServe(addr, router).Error())
}
