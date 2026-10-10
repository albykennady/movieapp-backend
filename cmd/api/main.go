package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"movieapp-backend/internal/repository"
	"movieapp-backend/internal/repository/dbrepo"
)

const port = 8080

type application struct {
	DSN          string
	Domain       string
	DB           repository.DataRepository
	auth         Auth
	JWTSecret    string
	JWTIssuer    string
	JWTAudience  string
	CookieDomain string
}

func main() {
	// set application config

	var app application

	// read from command line
	flag.StringVar(&app.DSN, "dsn", "host=localhost port=5432 user=postgres password=postgres dbname=movies sslmode=disable timezone=UTC connect_timeout=5", "Postgres connection string")
	flag.StringVar(&app.JWTSecret, "jwt-secret", "verysecret", "signing secret")
	flag.StringVar(&app.JWTIssuer, "jwt-issuer", "example.com", "signing issuer")
	flag.StringVar(&app.JWTAudience, "jwt-audience", "example.com", "JWT audience")
	flag.StringVar(&app.CookieDomain, "cookie-domain", "localhost", "Cookie domain")
	flag.StringVar(&app.Domain, "domain", "localhost", "Domain")
	flag.Parse()

	//connect to db
	conn, err := app.connectedToDB()
	if err != nil {
		log.Fatal(err)
	}
	app.DB = &dbrepo.PostgresDBRepo{
		DB: conn,
	}

	app.auth = Auth{
		Issuer:        app.JWTIssuer,
		Audience:      app.JWTAudience,
		Secret:        app.JWTSecret,
		TokenExpirey:  15 * time.Minute,
		RefreshExpiry: 24 * time.Hour,
		CookiePath:    "/",
		CookieName:    "__Host-refresh-token",
		CookieDomain:  app.CookieDomain,
	}

	defer app.DB.Connection().Close()

	fmt.Println("Connected to database")

	log.Println("Starting application on port", port)

	// start web server
	err = http.ListenAndServe(fmt.Sprintf(":%d", port), app.routes())
	if err != nil {
		log.Fatal(err)
	}

}
