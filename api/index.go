package handler

import (
	"log"
	"net/http"
	"sync"

	"agriculture-api/internal/config"
	"agriculture-api/internal/server"
)

var (
	initOnce sync.Once
	app      http.Handler
	initErr  error
)

func Handler(w http.ResponseWriter, r *http.Request) {
	initOnce.Do(func() {
		app, initErr = server.NewHandler(config.Load())
	})
	if initErr != nil {
		log.Printf("failed to initialize API handler: %v", initErr)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	app.ServeHTTP(w, r)
}
