package repository

import "github.com/zsbahtiar/lionparcel-test/internal/pkg/database"

type backOfficeRepository struct {
	db database.Postgres
}

type BackofficeRepository interface {
}

func NewBackOfficeRepository(db database.Postgres) BackofficeRepository {
	return &backOfficeRepository{
		db: db,
	}
}
