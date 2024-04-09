package handler

import (
	"encoding/json"
	"net/http"

	"github.com/carmooo/radarr-list/scraper"
	"github.com/go-chi/chi/v5"
)

func HandleLetterboxd(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "*")
	movies := scraper.GetMoviesFromLetterboxd(slug)

	jsonData, err := json.Marshal(movies)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	w.Write(jsonData)
}
