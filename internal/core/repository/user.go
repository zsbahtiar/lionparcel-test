package repository

import (
	"context"

	"github.com/zsbahtiar/lionparcel-test/internal/core/model/entity"
	"github.com/zsbahtiar/lionparcel-test/internal/pkg/database"
)

type userRepository struct {
	db database.Postgres
}

type UserRepository interface {
	CreateUser(ctx context.Context, user *entity.User) error
}

func NewUserRepository(db database.Postgres) UserRepository {
	return &userRepository{
		db: db,
	}
}

func (u *userRepository) CreateUser(ctx context.Context, user *entity.User) error {
	query := `
	INSERT INTO users (id, username, email, password_hash, is_admin) VALUES ($1, $2, $3, $4, $5)
	`
	_, err := u.db.Exec(ctx, query, user.ID, user.Username, user.Email, user.PasswordHash, user.IsAdmin)
	// @Todo handling custom error
	return err
}
