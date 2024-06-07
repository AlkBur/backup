package routes

import (
	"fmt"
	"log/slog"
	"net/http"

	"backup/log"
)

func NewRouter(logger *slog.Logger) http.Handler {
	// mux router
	mux := http.NewServeMux()

	// Routes
	mux.HandleFunc("GET /{$}", ShowIndexPage)

	http.HandleFunc("/api/greeting", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, HTMX from Go!")
	})

	mux.HandleFunc("GET /path/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "got path\n")
	})

	mux.HandleFunc("/task/{id}/", func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		fmt.Fprintf(w, "handling task with id=%v\n", id)
	})

	fs := http.FileServer(http.Dir("./static/"))
	mux.Handle("/static/", http.StripPrefix("/static", fs))

	// Middleware
	handler := log.Recovery(mux)
	handler = log.New(logger.WithGroup("http"))(handler)
	return handler
}
