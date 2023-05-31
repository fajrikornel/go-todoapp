package main

import (
	"fmt"
	"github.com/fajrikornel/go-todoapp/internal/config"
	"github.com/fajrikornel/go-todoapp/internal/db"
	"github.com/fajrikornel/go-todoapp/internal/server"
	"log"
	"net/http"
)

func main() {
	conf, err := config.GetConfig()
	if err != nil {
		log.Fatal("ERROR GETTING CONFIG: {}", err.Error())
		return
	}

	store, err := db.GetSqlStore(conf)
	if err != nil {
		log.Fatal("ERROR INITIALIZING DB: {}", err.Error())
		return
	}

	router := server.GetRouter(&store)

	addr := fmt.Sprintf(":%d", conf.AppPort)

	log.Println("STARTING WEB SERVER")
	log.Fatal(http.ListenAndServe(addr, router))
}
