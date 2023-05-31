package main

import (
	"fmt"
	"github.com/fajrikornel/go-todoapp/internal/config"
	"github.com/fajrikornel/go-todoapp/internal/server"
	"log"
	"net/http"
)

func main() {
	conf, err := config.GetConfig()
	if err != nil {
		log.Fatal("ERROR GETTING CONFIG: {}", err.Error())
	}

	router := server.GetRouter()

	addr := fmt.Sprintf(":%d", conf.AppPort)

	log.Println("STARTING WEB SERVER")
	log.Fatal(http.ListenAndServe(addr, router))
}
