package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"movieapp-backend/internal/repository"
	"movieapp-backend/internal/repository/dbrepo"
)

const port = 8080

type application struct {
	DSN    string
	Domain string
	DB     repository.DataRepository
}

func main() {
	// set application config

	var app application

	// read from command line
	flag.StringVar(&app.DSN, "dsn", "host=localhost port=5432 user=postgres password=postgres dbname=movies sslmode=disable timezone=UTC connect_timeout=5", "Postgres connection string")
	flag.Parse()

	//connect to db
	conn, err := app.connectedToDB()
	if err != nil {
		log.Fatal(err)
	}
	app.DB = &dbrepo.PostgresDBRepo{
		DB: conn,
	}

	defer app.DB.Connection().Close()

	fmt.Println("Connected to database")

	app.Domain = "example.com"

	log.Println("Starting application on port", port)

	// start web server
	err = http.ListenAndServe(fmt.Sprintf(":%d", port), app.routes())
	if err != nil {
		log.Fatal(err)
	}

}
