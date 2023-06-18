package main

import (
	"context"
	"github.com/fajrikornel/go-todoapp/internal/config"
	"github.com/fajrikornel/go-todoapp/internal/db"
	"github.com/fajrikornel/go-todoapp/internal/logging"
)

func main() {
	dbConfig := config.GetTestDbConfig()
	sqlStore, err := db.GetSqlStore(&dbConfig)
	if err != nil {
		logging.Errorf(context.Background(), "ERROR INITIALIZING TEST DB: %v\n", err.Error())
		return
	}

	err = sqlStore.DoMigrations()
	if err != nil {
		logging.Errorf(context.Background(), "ERROR EXECUTING TEST MIGRATIONS: %v\n", err.Error())
		return
	}

	logging.Infof(context.Background(), "SUCCESSFUL EXECUTING TEST MIGRATIONS")
}
