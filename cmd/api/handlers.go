package main

import (
	"encoding/json"
	"fmt"
	"movieapp-backend/internal/models"
	"net/http"
	"time"
)

func (app *application) Home(w http.ResponseWriter, r *http.Request) {
	var payload = struct {
		Status  string `json:"status"`
		Message string `json:"message"`
		Version string `json:"version"`
	}{

		Status:  "success",
		Message: "Go Movie API is up and running",
		Version: "1.0.0",
	}
	out, err := json.Marshal(payload)
	if err != nil {
		fmt.Println(err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(out)
}

func (app *application) AllMovies(w http.ResponseWriter, r *http.Request) {
	var movies []models.Movie

	highlander := models.Movie{
		ID:          1,
		Title:       "Highlander",
		Year:        1986,
		Runtime:     116,
		MPAARating:  "R",
		Description: "An immortal warrior battles through centuries of conflict.",
		CreatedAt:   time.Now().Format(time.RFC3339),
		UpdatedAt:   time.Now().Format(time.RFC3339),
	}

	movies = append(movies, highlander)

	inception := models.Movie{
		ID:          2,
		Title:       "Inception",
		Year:        2010,
		Runtime:     148,
		MPAARating:  "PG-13",
		Description: "A thief enters people's dreams to steal and plant information.",
		CreatedAt:   time.Now().Format(time.RFC3339),
		UpdatedAt:   time.Now().Format(time.RFC3339),
	}

	movies = append(movies, inception)

	out, err := json.Marshal(movies)
	if err != nil {
		fmt.Println(err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(out)

}
