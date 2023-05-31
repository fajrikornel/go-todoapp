package main

import (
	"fmt"
	"github.com/fajrikornel/go-todoapp/internal/config"
	"github.com/fajrikornel/go-todoapp/internal/db"
)

func main() {
	conf, err := config.GetConfig()
	if err != nil {
		fmt.Printf("ERROR GETTING CONFIG: %v\n", err.Error())
		return
	}

	sqlStore, err := db.GetSqlStore(conf)
	if err != nil {
		fmt.Printf("ERROR INITIALIZING DB: %v\n", err.Error())
		return
	}

	err = sqlStore.DoMigrations()
	if err != nil {
		fmt.Printf("ERROR EXECUTING MIGRATIONS: %v\n", err.Error())
		return
	}

	fmt.Println("SUCCESSFUL EXECUTING MIGRATIONS")
}
