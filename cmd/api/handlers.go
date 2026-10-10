package main

import (
	"log"
	"net/http"
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
