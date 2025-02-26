package repository

import (
	"context"

	"github.com/zsbahtiar/lionparcel-test/internal/core/model/entity"
	"github.com/zsbahtiar/lionparcel-test/internal/pkg/database"
)

type backOfficeRepository struct {
	db database.Postgres
}

type BackofficeRepository interface {
	CreateMovice(ctx context.Context, movie *entity.Movie) error
	UpdateMovice(ctx context.Context, movie *entity.Movie) error
}

func NewBackOfficeRepository(db database.Postgres) BackofficeRepository {
	return &backOfficeRepository{
		db: db,
	}
}

func (b *backOfficeRepository) CreateMovice(ctx context.Context, movie *entity.Movie) error {
	query := `
	INSERT INTO movies (id, title, description, duration, link, genres, artists)
	`
	_, err := b.db.Exec(ctx, query, movie.Title, movie.Description, movie.Duration, movie.Link, movie.Genres, movie.Artists)
	// @Todo handling custom error
	return err
}

func (b *backOfficeRepository) UpdateMovice(ctx context.Context, movie *entity.Movie) error {
	query := `
	UPDATE movies
	SET title = $1, description = $2, duration = $3, link = $4, genres = $5, artists = $6
	WHERE id = $7
	`

	_, err := b.db.Exec(ctx, query, movie.Title, movie.Description, movie.Duration, movie.Link, movie.Genres, movie.Artists, movie.ID)
	// @Todo handling custom error
	return err
}
