package repository

import (
	"context"
	"errors"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/zsbahtiar/lionparcel-test/internal/core/model/entity"
	"github.com/zsbahtiar/lionparcel-test/internal/core/model/internalerror"
	"github.com/zsbahtiar/lionparcel-test/internal/core/model/request"
	"github.com/zsbahtiar/lionparcel-test/internal/pkg/database"
)

type movieRepository struct {
	db database.Postgres
}

type MovieRepository interface {
	CreateMovie(ctx context.Context, movie *entity.Movie) error
	UpdateMovie(ctx context.Context, movie *entity.Movie) error
	GetMovies(ctx context.Context, req *request.GetMovies) ([]entity.Movie, int64, error)
	GetMovie(ctx context.Context, movieID string) (*entity.Movie, error)
	GetViewMovies(ctx context.Context, movieID string) (int64, error)
	GetMostViewedMovies(ctx context.Context, limit int) ([]entity.MovieViewCount, error)
	GetMostViewedGenres(ctx context.Context, limit int) ([]entity.GenreViewCount, error)
	GetMostVotedMovies(ctx context.Context) ([]entity.MovieVotedCount, error)
	CreateVote(ctx context.Context, vote *entity.Vote) error
	DeleteVote(ctx context.Context, userID, movieID string) error
	GetVotedMovieOfUser(ctx context.Context, userID string) ([]entity.UserMovieVote, error)
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
	return err
}

func (m *movieRepository) UpdateMovie(ctx context.Context, movie *entity.Movie) error {
	query := `
	UPDATE movies
	SET title = $1, description = $2, duration = $3, link = $4, genres = $5, artists = $6
	WHERE id = $7
	`

	res, err := m.db.Exec(ctx, query, movie.Title, movie.Description, movie.Duration, movie.Link, movie.Genres, movie.Artists, movie.ID)
	if err != nil {
		return err
	}

	if res.RowsAffected() < 1 {
		return internalerror.ErrMovieNotFound
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
		if err == pgx.ErrNoRows {
			return nil, internalerror.ErrMovieNotFound
		}
		return nil, err
	}
	return &movie, nil
}

func (m *movieRepository) GetViewMovies(ctx context.Context, movieID string) (int64, error) {
	query := `
SELECT COUNT(*) AS total_views
FROM user_movie_views
WHERE movie_id = $1
`
	var totalViews int64
	err := m.db.SelectOne(ctx, &totalViews, query, movieID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return 0, internalerror.ErrMovieNotFound
		}
		return 0, err
	}
	return totalViews, nil
}

func (m *movieRepository) GetMostViewedMovies(ctx context.Context, limit int) ([]entity.MovieViewCount, error) {
	query := `
    SELECT 
        m.id,
        m.title,
        COUNT(umv.id) AS view_count
    FROM 
        movies m
    JOIN 
        user_movie_views umv ON m.id = umv.movie_id
    GROUP BY 
        m.id, m.title
    ORDER BY 
        view_count DESC
    LIMIT $1
    `
	var movies []entity.MovieViewCount
	err := m.db.Select(ctx, &movies, query, limit)
	if err != nil {
		return nil, err
	}
	return movies, nil
}

func (m *movieRepository) GetMostViewedGenres(ctx context.Context, limit int) ([]entity.GenreViewCount, error) {
	query := `
    SELECT 
        genre,
        COUNT(*) AS view_count
    FROM 
        user_movie_views umv
    JOIN 
        movies m ON umv.movie_id = m.id,
        jsonb_array_elements_text(m.genres) AS genre
    GROUP BY 
        genre
    ORDER BY 
        view_count DESC
    LIMIT $1
    `
	var genres []entity.GenreViewCount
	err := m.db.Select(ctx, &genres, query, limit)
	if err != nil {
		return nil, err
	}
	return genres, nil
}

func (m *movieRepository) GetMostVotedMovies(ctx context.Context) ([]entity.MovieVotedCount, error) {
	query := `
	SELECT 
		m.id,
		m.title,
		COUNT(v.id) AS voted_count
	FROM 
		movies m
	JOIN 
		votes v ON m.id = v.movie_id
	GROUP BY 
		m.id, m.title
	ORDER BY 
		voted_count DESC
	`
	var movies []entity.MovieVotedCount
	err := m.db.Select(ctx, &movies, query)
	if err != nil {
		return nil, err
	}
	return movies, nil
}

func (m *movieRepository) CreateVote(ctx context.Context, vote *entity.Vote) error {
	query := `
	INSERT INTO votes (user_id, movie_id) VALUES ($1, $2)
	`
	_, err := m.db.Exec(ctx, query, vote.UserID, vote.MovieID)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				return internalerror.ErrVoteExist
			}
		}
		return err
	}
	return nil
}

func (m *movieRepository) DeleteVote(ctx context.Context, userID, movieID string) error {
	query := `
	DELETE FROM votes
	WHERE user_id = $1 AND movie_id = $2
	`
	_, err := m.db.Exec(ctx, query, userID, movieID)
	if err != nil {
		return err
	}
	return nil
}

func (m *movieRepository) GetVotedMovieOfUser(ctx context.Context, userID string) ([]entity.UserMovieVote, error) {
	query := `
	SELECT 
		m.id,
		m.title,
		v.created_at AS voted_at
	FROM 
		movies m
	JOIN 
		votes v ON m.id = v.movie_id
	WHERE 
		v.user_id = $1
	ORDER BY v.created_at DESC
	`
	var votes []entity.UserMovieVote
	err := m.db.Select(ctx, &votes, query, userID)
	if err != nil {
		return nil, err
	}
	return votes, nil
}
