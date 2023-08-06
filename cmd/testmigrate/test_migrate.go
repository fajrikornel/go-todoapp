package main

import (
	"context"
	"fmt"
	"github.com/fajrikornel/go-todoapp/internal/config"
	"github.com/fajrikornel/go-todoapp/internal/logging"
	"github.com/golang-migrate/migrate"
	_ "github.com/golang-migrate/migrate/database/postgres"
	_ "github.com/golang-migrate/migrate/source/file"
	"path/filepath"
	"runtime"
)

func main() {
	dbConfig := config.GetTestDbConfig()

	m, err := migrate.New(
		getMigrationsPath(),
		getDbUrl(dbConfig),
	)

	if err != nil {
		panic(err)
	}

	err = m.Up()
	if err != nil {
		panic(err)
	}

	logging.Infof(context.Background(), "SUCCESSFUL EXECUTING MIGRATIONS")
}

func getMigrationsPath() string {
	_, b, _, _ := runtime.Caller(0)
	path := filepath.Dir(b) + "/../../migrations"

	return "file://" + path
}

func getDbUrl(dbConfig config.DbConfig) string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=disable", // todo: enable SSL for secure connections
		dbConfig.DbUsername,
		dbConfig.DbPassword,
		dbConfig.DbHost,
		dbConfig.DbPort,
		dbConfig.DbName,
	)
}
