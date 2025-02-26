package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/zsbahtiar/lionparcel-test/internal/core/model/entity"
	"github.com/zsbahtiar/lionparcel-test/internal/core/model/request"
	"github.com/zsbahtiar/lionparcel-test/internal/pkg/database"
)

type backOfficeRepository struct {
	db database.Postgres
}

type BackofficeRepository interface {
	CreateMovie(ctx context.Context, movie *entity.Movie) error
	UpdateMovie(ctx context.Context, movie *entity.Movie) error
	// @Todo: add more request and response
	GetMovies(ctx context.Context, req *request.GetMovies) ([]entity.Movie, error)
}

func NewBackOfficeRepository(db database.Postgres) BackofficeRepository {
	return &backOfficeRepository{
		db: db,
	}
}

func (b *backOfficeRepository) CreateMovie(ctx context.Context, movie *entity.Movie) error {
	query := `
	INSERT INTO movies (id, title, description, duration, link, genres, artists) VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	_, err := b.db.Exec(ctx, query, movie.ID, movie.Title, movie.Description, movie.Duration, movie.Link, movie.Genres, movie.Artists)
	// @Todo handling custom error
	return err
}

func (b *backOfficeRepository) UpdateMovie(ctx context.Context, movie *entity.Movie) error {
	query := `
	UPDATE movies
	SET title = $1, description = $2, duration = $3, link = $4, genres = $5, artists = $6
	WHERE id = $7
	`

	res, err := b.db.Exec(ctx, query, movie.Title, movie.Description, movie.Duration, movie.Link, movie.Genres, movie.Artists, movie.ID)
	// @Todo handling custom error
	if err != nil {
		return err
	}

	if res.RowsAffected() < 1 {
		return fmt.Errorf("movie with id %s not found", movie.ID)
	}
	return nil
}

func (b *backOfficeRepository) GetMovies(ctx context.Context, req *request.GetMovies) ([]entity.Movie, error) {
	query := `
SELECT id, title, description, duration, link, genres, artists
FROM movies
`
	var conditions []string
	var values []interface{}

	if len(req.Search) > 0 {
		// using % wildcard in last, so it will search for any string that start with search value, cause performance issue // im not create full text search index
		searchValue := req.Search + "%"
		searchConditions := []string{
			"title ILIKE ?",
			"description ILIKE ?",
			"array_to_string(artists, ',') ILIKE ?",
			"array_to_string(genres, ',') ILIKE ?",
		}
		conditions = append(conditions, "("+strings.Join(searchConditions, " OR ")+")")
		values = append(values, searchValue)
	}

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	query += " ORDER BY title ASC"

	if req.Limit > 0 {
		query += " LIMIT ?"
		values = append(values, req.Limit)

		if req.Page > 0 {
			offset := (req.Page - 1) * req.Limit
			query += " OFFSET ?"
			values = append(values, offset)
		}
	}

	var movies []entity.Movie
	err := b.db.Select(ctx, &movies, b.db.Rebind(query), values...)
	// @Todo handling custom error
	return movies, err
}
