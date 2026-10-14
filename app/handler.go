package app

import (
	"net/http"

	"agriculture-api/internal/config"
	"agriculture-api/internal/server"
)

func NewHandler() (http.Handler, error) {
	return server.NewHandler(config.Load())
}
