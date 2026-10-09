package repository

import (
	"database/sql"
	"movieapp-backend/internal/models"
)

type DataRepository interface {
	Connection() *sql.DB
	AllMovies() ([]*models.Movie, error)
}
