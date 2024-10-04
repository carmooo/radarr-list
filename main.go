package main

import (
	"net/http"

	"github.com/carmooo/radarr-list/handler"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/letterboxd/*", handler.HandleLetterboxd)
	http.ListenAndServe(":3000", r)
}
