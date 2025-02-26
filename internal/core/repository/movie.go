package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/zsbahtiar/lionparcel-test/internal/core/model/entity"
	"github.com/zsbahtiar/lionparcel-test/internal/core/model/request"
	"github.com/zsbahtiar/lionparcel-test/internal/pkg/database"
)

type movieRepository struct {
	db database.Postgres
}

type MovieRepository interface {
	CreateMovie(ctx context.Context, movie *entity.Movie) error
	UpdateMovie(ctx context.Context, movie *entity.Movie) error
	// @Todo: add more request and response
	GetMovies(ctx context.Context, req *request.GetMovies) ([]entity.Movie, int64, error)
	GetMovie(ctx context.Context, movieID string) (*entity.Movie, error)
}

func NewMovieRepository(db database.Postgres) MovieRepository {
	return &movieRepository{
		db: db,
	}
}

func (m *movieRepository) CreateMovie(ctx context.Context, movie *entity.Movie) error {
	query := `
	INSERT INTO movies (id, title, description, duration, link, genres, artists) VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	_, err := m.db.Exec(ctx, query, movie.ID, movie.Title, movie.Description, movie.Duration, movie.Link, movie.Genres, movie.Artists)
	// @Todo handling custom error
	return err
}

func (m *movieRepository) UpdateMovie(ctx context.Context, movie *entity.Movie) error {
	query := `
	UPDATE movies
	SET title = $1, description = $2, duration = $3, link = $4, genres = $5, artists = $6
	WHERE id = $7
	`

	res, err := m.db.Exec(ctx, query, movie.Title, movie.Description, movie.Duration, movie.Link, movie.Genres, movie.Artists, movie.ID)
	// @Todo handling custom error
	if err != nil {
		return err
	}

	if res.RowsAffected() < 1 {
		return fmt.Errorf("movie with id %s not found", movie.ID)
	}
	return nil
}

func (m *movieRepository) GetMovies(ctx context.Context, req *request.GetMovies) ([]entity.Movie, int64, error) {
	query := `
SELECT id, title, description, duration, link, genres, artists, COUNT(*) OVER() AS total_count
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
			"EXISTS (SELECT 1 FROM jsonb_array_elements_text(artists) AS a WHERE a ILIKE ?)",
			"EXISTS (SELECT 1 FROM jsonb_array_elements_text(genres) AS g WHERE g ILIKE ?)",
		}
		conditions = append(conditions, "("+strings.Join(searchConditions, " OR ")+")")
		values = append(values, searchValue, searchValue, searchValue, searchValue)
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

	type movieWithCount struct {
		entity.Movie
		TotalCount int64 `db:"total_count"`
	}
	var movies []entity.Movie
	var totalCount int64 = 0
	var moviesWithCount []movieWithCount
	err := m.db.Select(ctx, &moviesWithCount, m.db.Rebind(query), values...)
	if err != nil {
		// @Todo handling custom error
		return nil, 0, err
	}
	movies = make([]entity.Movie, len(moviesWithCount))
	for i, movie := range moviesWithCount {
		movies[i] = movie.Movie
		totalCount = movie.TotalCount
	}
	return movies, totalCount, nil
}

func (m *movieRepository) GetMovie(ctx context.Context, movieID string) (*entity.Movie, error) {
	query := `
SELECT id, title, description, duration, link, genres, artists
FROM movies
WHERE id = $1
`
	var movie entity.Movie
	err := m.db.SelectOne(ctx, &movie, query, movieID)
	if err != nil {
		// @Todo handling custom error
		return nil, err
	}
	return &movie, nil
}
