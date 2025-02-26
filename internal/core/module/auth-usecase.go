package module

import (
	"context"

	"github.com/google/uuid"
	"github.com/zsbahtiar/lionparcel-test/internal/core/model/entity"
	"github.com/zsbahtiar/lionparcel-test/internal/core/model/request"
	"github.com/zsbahtiar/lionparcel-test/internal/core/repository"
	"golang.org/x/crypto/bcrypt"
)

type authUsecase struct {
	userRepository repository.UserRepository
}

type AuthUsecase interface {
	RegisterUser(ctx context.Context, req *request.Register) error
}

func NewAuthUsecase(userRepository repository.UserRepository) AuthUsecase {
	return &authUsecase{userRepository: userRepository}
}

func (a *authUsecase) RegisterUser(ctx context.Context, req *request.Register) error {
	paswordHash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user := &entity.User{
		ID:           uuid.New().String(),
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: string(paswordHash),
		IsAdmin:      req.IsAdmin,
	}
	return a.userRepository.CreateUser(ctx, user)
}
