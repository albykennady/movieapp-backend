package main

import (
	"database/sql"
	"errors"
	"log"
	"movieapp-backend/internal/models"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
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
	_ = app.writeJSON(w, http.StatusOK, payload)
}

func (app *application) AllMovies(w http.ResponseWriter, r *http.Request) {
	movies, err := app.DB.AllMovies()
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	_ = app.writeJSON(w, http.StatusOK, movies)
}

func (app *application) authenticate(w http.ResponseWriter, r *http.Request) {
	//read JSON payload
	//validate user gainst DB
	//check passsword
	//creare JWT user

	u := &jwtUser{
		ID:        1,
		FirstName: "Admin",
		LastName:  "User",
	}

	//genereate tokens
	tokens, err := app.auth.GenerateTokenPair(u)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	log.Println("Generated tokens:", tokens.Token)

	w.Write([]byte(tokens.Token))
}

func (app *application) GetMovie(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	movieID, err := strconv.Atoi(id)

	if err != nil {
		app.errorJSON(w, err)
		return
	}

	movie, err := app.DB.OneMovie(movieID)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	_ = app.writeJSON(w, http.StatusOK, movie)
}

func (app *application) MovieForEdit(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	movieID, err := strconv.Atoi(id)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	movie, genres, err := app.DB.OneMovieForEdit(movieID)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	var payload = struct {
		Movie  *models.Movie   `json: "movie"`
		Genres []*models.Genre `json: "genres"`
	}{
		movie,
		genres,
	}

	_ = app.writeJSON(w, http.StatusOK, payload)

}

func (app *application) CreateMovie(w http.ResponseWriter, r *http.Request) {
	var payload models.Movie
	if err := app.readJSON(w, r, &payload); err != nil {
		app.errorJSON(w, err)
		return
	}

	if strings.TrimSpace(payload.Title) == "" {
		app.errorJSON(w, errors.New("movie title is required"))
		return
	}

	if payload.Runtime < 0 {
		app.errorJSON(w, errors.New("runtime must be zero or greater"))
		return
	}

	movieID, err := app.DB.InsertMovie(&payload)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	response := map[string]interface{}{
		"id":      movieID,
		"message": "movie created successfully",
	}

	_ = app.writeJSON(w, http.StatusCreated, response)
}

func (app *application) UpdateMovie(w http.ResponseWriter, r *http.Request) {
	movieID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil || movieID < 1 {
		app.errorJSON(w, errors.New("movie id must be a positive integer"), http.StatusBadRequest)
		return
	}

	var payload models.MovieUpdate
	if err := app.readJSON(w, r, &payload); err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}

	if err := app.DB.UpdateMovie(movieID, &payload); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			app.errorJSON(w, err, http.StatusNotFound)
			return
		}
		app.errorJSON(w, err)
		return
	}

	_ = app.writeJSON(w, http.StatusOK, map[string]interface{}{
		"id":      movieID,
		"message": "movie updated successfully",
	})
}

func (app *application) DeleteMovie(w http.ResponseWriter, r *http.Request) {
	movieID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil || movieID < 1 {
		app.errorJSON(w, errors.New("movie id must be a positive integer"), http.StatusBadRequest)
		return
	}

	if err := app.DB.DeleteMovie(movieID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			app.errorJSON(w, err, http.StatusNotFound)
			return
		}
		app.errorJSON(w, err)
		return
	}

	_ = app.writeJSON(w, http.StatusOK, map[string]interface{}{
		"id":      movieID,
		"message": "movie deleted successfully",
	})
}
