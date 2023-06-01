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
	conf, err := config.GetConfig()
	if err != nil {
		logging.Errorf(context.Background(), "ERROR GETTING CONFIG: {}", err.Error())
		return
	}

	store, err := db.GetSqlStore(conf)
	if err != nil {
		logging.Errorf(context.Background(), "ERROR INITIALIZING DB: {}", err.Error())
		return
	}

	projectRepository := repository.NewProjectRepository(store)
	itemRepository := repository.NewItemRepository(store)
	router := server.GetRouter(projectRepository, itemRepository)

	addr := fmt.Sprintf(":%d", conf.AppPort)

	logging.Infof(context.Background(), "STARTING WEB SERVER")
	logging.Errorf(context.Background(), "%s", http.ListenAndServe(addr, router).Error())
}
