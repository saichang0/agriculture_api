package main

import (
	"log"
	"net/http"

	"agriculture-api/internal/config"
	"agriculture-api/internal/server"
)

func main() {
	cfg := config.Load()
	app, err := server.NewHandler(cfg)
	if err != nil {
		log.Fatalf("failed to initialize API server: %v", err)
	}

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", cfg.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.Port, app))
}
