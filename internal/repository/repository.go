package repository

import (
	"database/sql"
	"movieapp-backend/internal/models"
)

type DataRepository interface {
	Connection() *sql.DB
	AllMovies() ([]*models.Movie, error)
	OneMovieForEdit(id int) (*models.Movie, []*models.Genre, error)
	OneMovie(id int) (*models.Movie, error)
	InsertMovie(movie *models.Movie) (int, error)
	UpdateMovie(id int, movie *models.MovieUpdate) error
	DeleteMovie(id int) error
}
