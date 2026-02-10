package main

import (
	"log/slog"
	"net/http"

	"github.com/j4y_funabashi/inari/apps/api/pkg/webhandler"
)

func main() {
	logger := slog.Default()
	router := webhandler.NewWebHandler()
	port := ":8080"
	logger.Info("inari web server running on port", "port", port)
	http.ListenAndServe(port, router)
}
