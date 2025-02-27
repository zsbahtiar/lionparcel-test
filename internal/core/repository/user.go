package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/zsbahtiar/lionparcel-test/internal/core/model/entity"
	"github.com/zsbahtiar/lionparcel-test/internal/core/model/internalerror"
	"github.com/zsbahtiar/lionparcel-test/internal/pkg/database"
)

type userRepository struct {
	db database.Postgres
}

type UserRepository interface {
	CreateUser(ctx context.Context, user *entity.User) error
	GetUserByEmail(ctx context.Context, email string) (*entity.User, error)
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
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				return internalerror.ErrUserExist
			}
		}
		return err
	}
	return nil
}

func (u *userRepository) GetUserByEmail(ctx context.Context, email string) (*entity.User, error) {
	query := `
	SELECT id, username, email, password_hash, is_admin
	FROM users
	WHERE email = $1 LIMIT 1
	`
	user := &entity.User{}
	err := u.db.SelectOne(ctx, user, query, email)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, internalerror.ErrUserNotFound
		}
		return nil, err
	}
	return user, nil
}
