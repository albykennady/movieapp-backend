package dbrepo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"movieapp-backend/internal/models"
	"strings"
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
		var year sql.NullInt64

		err := rows.Scan(
			&movie.ID,
			&movie.Title,
			&movie.Description,
			&year,
			&movie.Runtime,
			&movie.MPAARating,
			&movie.Image,
			&movie.CreatedAt,
			&movie.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		if year.Valid {
			movie.Year = int(year.Int64)
		}

		movies = append(movies, &movie)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return movies, nil
}

func (m *PostgresDBRepo) OneMovieForEdit(id int) (*models.Movie, []*models.Genre, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `select id, title, runtime, mpaa_rating,
	description, coalesce(image, ''), created_at, updated_at
	from movies where id = $1`

	row := m.DB.QueryRowContext(ctx, query, id)

	var movie models.Movie

	err := row.Scan(
		&movie.ID,
		&movie.Title,
		&movie.Runtime,
		&movie.MPAARating,
		&movie.Description,
		&movie.Image,
		&movie.CreatedAt,
		&movie.UpdatedAt,
	)

	if err != nil {
		return nil, nil, err
	}

	//get genres if any
	query = `select g.id, g.genre from movies_genres mg
	left join genres g on (mg.genre_id = g.id)
	where mg.movie_id = $1
	order by g.genre`

	rows, err := m.DB.QueryContext(ctx, query, id)
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()

	var genres []*models.Genre
	var genresArrray []int

	for rows.Next() {
		var g models.Genre
		err := rows.Scan(
			&g.ID,
			&g.Genre,
		)
		if err != nil {
			return nil, nil, err
		}

		genres = append(genres, &g)
		genresArrray = append(genresArrray, g.ID)
	}

	movie.Genres = genres
	movie.GenresArray = genresArrray

	var allGenres []*models.Genre

	query = `select id, genre from genres order by genre`
	gRows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, nil, err
	}
	defer gRows.Close()

	for gRows.Next() {
		var g models.Genre
		err := gRows.Scan(
			&g.ID,
			&g.Genre,
		)

		if err != nil {
			return nil, nil, err
		}

		allGenres = append(allGenres, &g)

	}
	return &movie, allGenres, err
}

func (m *PostgresDBRepo) OneMovie(id int) (*models.Movie, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `select id, title, EXTRACT(YEAR FROM release_date)::int,
	runtime, mpaa_rating, description, coalesce(image, ''),
	created_at, updated_at
	from movies where id = $1`

	row := m.DB.QueryRowContext(ctx, query, id)

	var movie models.Movie
	var year sql.NullInt64

	err := row.Scan(
		&movie.ID,
		&movie.Title,
		&year,
		&movie.Runtime,
		&movie.MPAARating,
		&movie.Description,
		&movie.Image,
		&movie.CreatedAt,
		&movie.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}
	if year.Valid {
		movie.Year = int(year.Int64)
	}

	//get genres if any
	query = `select g.id, g.genre from movies_genres mg
	left join genres g on (mg.genre_id = g.id)
	where mg.movie_id = $1
	order by g.genre`

	rows, err := m.DB.QueryContext(ctx, query, id)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	defer rows.Close()

	var genres []*models.Genre
	for rows.Next() {
		var g models.Genre
		err := rows.Scan(
			&g.ID,
			&g.Genre,
		)
		if err != nil {
			return nil, err
		}

		genres = append(genres, &g)
	}

	movie.Genres = genres

	return &movie, err
}

func (m *PostgresDBRepo) InsertMovie(movie *models.Movie) (int, error) {
	if movie == nil {
		return 0, errors.New("movie is required")
	}

	if strings.TrimSpace(movie.Title) == "" {
		return 0, errors.New("movie title is required")
	}

	if movie.Runtime < 0 {
		return 0, errors.New("runtime must be zero or greater")
	}

	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	tx, err := m.DB.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	var releaseDate sql.NullString
	if movie.Year > 0 {
		releaseDate = sql.NullString{
			String: fmt.Sprintf("%d-01-01", movie.Year),
			Valid:  true,
		}
	}

	var movieID int
	query := `
		INSERT INTO movies (title, release_date, runtime, mpaa_rating, description, image, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, NOW(), NOW())
		RETURNING id
	`

	err = tx.QueryRowContext(ctx, query,
		strings.TrimSpace(movie.Title),
		releaseDate,
		movie.Runtime,
		strings.TrimSpace(movie.MPAARating),
		movie.Description,
		movie.Image,
	).Scan(&movieID)
	if err != nil {
		return 0, err
	}

	for _, genreID := range movie.GenresArray {
		_, err = tx.ExecContext(ctx, `
			INSERT INTO movies_genres (movie_id, genre_id)
			VALUES ($1, $2)
		`, movieID, genreID)
		if err != nil {
			return 0, err
		}
	}

	if err = tx.Commit(); err != nil {
		return 0, err
	}

	return movieID, nil
}

func (m *PostgresDBRepo) UpdateMovie(id int, movie *models.MovieUpdate) error {
	if movie == nil {
		return errors.New("movie update is required")
	}

	updates := make([]string, 0, 7)
	args := make([]interface{}, 0, 8)
	addUpdate := func(column string, value interface{}) {
		args = append(args, value)
		updates = append(updates, fmt.Sprintf("%s = $%d", column, len(args)))
	}

	if movie.Title != nil {
		title := strings.TrimSpace(*movie.Title)
		if title == "" {
			return errors.New("movie title cannot be empty")
		}
		addUpdate("title", title)
	}
	if movie.Year != nil {
		if *movie.Year < 0 {
			return errors.New("year must be zero or greater")
		}
		if *movie.Year == 0 {
			addUpdate("release_date", nil)
		} else {
			addUpdate("release_date", fmt.Sprintf("%04d-01-01", *movie.Year))
		}
	}
	if movie.Runtime != nil {
		if *movie.Runtime < 0 {
			return errors.New("runtime must be zero or greater")
		}
		addUpdate("runtime", *movie.Runtime)
	}
	if movie.MPAARating != nil {
		addUpdate("mpaa_rating", strings.TrimSpace(*movie.MPAARating))
	}
	if movie.Description != nil {
		addUpdate("description", *movie.Description)
	}
	if movie.Image != nil {
		addUpdate("image", *movie.Image)
	}
	if len(updates) == 0 && movie.GenresArray == nil {
		return errors.New("at least one movie field must be provided")
	}

	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	tx, err := m.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if len(updates) > 0 {
		args = append(args, id)
		query := fmt.Sprintf(
			"UPDATE movies SET %s, updated_at = NOW() WHERE id = $%d",
			strings.Join(updates, ", "),
			len(args),
		)
		result, err := tx.ExecContext(ctx, query, args...)
		if err != nil {
			return err
		}
		rowsAffected, err := result.RowsAffected()
		if err != nil {
			return err
		}
		if rowsAffected == 0 {
			return sql.ErrNoRows
		}
	}

	if movie.GenresArray != nil {
		result, err := tx.ExecContext(ctx, "UPDATE movies SET updated_at = NOW() WHERE id = $1", id)
		if err != nil {
			return err
		}
		rowsAffected, err := result.RowsAffected()
		if err != nil {
			return err
		}
		if rowsAffected == 0 {
			return sql.ErrNoRows
		}

		if _, err = tx.ExecContext(ctx, "DELETE FROM movies_genres WHERE movie_id = $1", id); err != nil {
			return err
		}
		for _, genreID := range *movie.GenresArray {
			if _, err = tx.ExecContext(ctx, `
				INSERT INTO movies_genres (movie_id, genre_id)
				VALUES ($1, $2)
			`, id, genreID); err != nil {
				return err
			}
		}
	}

	return tx.Commit()
}

func (m *PostgresDBRepo) DeleteMovie(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	result, err := m.DB.ExecContext(ctx, "DELETE FROM movies WHERE id = $1", id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}
