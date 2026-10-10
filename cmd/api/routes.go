package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (app *application) routes() http.Handler {
	mux := chi.NewRouter()

	mux.Use(middleware.Recoverer)

	mux.Use(app.enableCORS)

	mux.Get("/", app.Home)

	mux.Get("/authenticate", app.authenticate)

	mux.Get("/movies/{id}", app.GetMovie)
	mux.Post("/movies", app.CreateMovie)
	mux.Patch("/movies/{id}", app.UpdateMovie)
	mux.Delete("/movies/{id}", app.DeleteMovie)
	mux.Get("/movies", app.AllMovies)

	return mux
}
