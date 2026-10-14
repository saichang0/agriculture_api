package handler

import (
	"log"
	"net/http"
	"sync"

	application "agriculture-api/app"
)

var (
	initOnce sync.Once
	app      http.Handler
	initErr  error
)

func Handler(w http.ResponseWriter, r *http.Request) {
	initOnce.Do(func() {
		app, initErr = application.NewHandler()
	})
	if initErr != nil {
		log.Printf("failed to initialize API handler: %v", initErr)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	app.ServeHTTP(w, r)
}
