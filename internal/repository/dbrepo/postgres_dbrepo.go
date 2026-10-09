package dbrepo

import (
	"context"
	"database/sql"
	"movieapp-backend/internal/models"
	"time"
)

type PostgresDBRepo struct {
	DB *sql.DB
}

const dbTimeout = time.Second * 3

func (m *PostgresDBRepo) Connection() *sql.DB {
	return m.DB
}

func (m *PostgresDBRepo) AllMovies() ([]*models.Movie, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `
        SELECT
        id,
        title,
        description,
        EXTRACT(YEAR FROM release_date)::int,
        runtime,
        mpaa_rating,
        COALESCE(image, ''),
        created_at,
        updated_at
    FROM movies
    ORDER BY title
    `

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	movies := make([]*models.Movie, 0)

	for rows.Next() {
		var movie models.Movie

		err := rows.Scan(
			&movie.ID,
			&movie.Title,
			&movie.Description,
			&movie.Year,
			&movie.Runtime,
			&movie.MPAARating,
			&movie.Image,
			&movie.CreatedAt,
			&movie.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		movies = append(movies, &movie)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return movies, nil
}
