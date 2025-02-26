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
