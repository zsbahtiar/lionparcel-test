package module

import (
	"context"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/zsbahtiar/lionparcel-test/internal/core/model/entity"
	"github.com/zsbahtiar/lionparcel-test/internal/core/model/request"
	"github.com/zsbahtiar/lionparcel-test/internal/core/model/response"
	"github.com/zsbahtiar/lionparcel-test/internal/core/repository"
	"golang.org/x/crypto/bcrypt"
)

type authUsecase struct {
	userRepository repository.UserRepository
	jwtSecret      string
}

type AuthUsecase interface {
	RegisterUser(ctx context.Context, req *request.Register) error
	Login(ctx context.Context, req *request.Login) (*response.Login, error)
}

func NewAuthUsecase(userRepository repository.UserRepository, jwtSecret string) AuthUsecase {
	return &authUsecase{userRepository: userRepository, jwtSecret: jwtSecret}
}

const (
	// @Todo: move to env
	AccessTokenExpiration = 1 * time.Hour
)

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

func (a *authUsecase) Login(ctx context.Context, req *request.Login) (*response.Login, error) {
	user, err := a.userRepository.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))
	if err != nil {
		return nil, err
	}

	now := time.Now()
	expireAt := now.Add(AccessTokenExpiration)
	// @Todo: add expiration time if still have time for development
	claims := jwt.MapClaims{
		"user_id": user.ID,
		"exp":     expireAt.Unix(),
		"iat":     now.Unix(),
		"type":    "access",
	}
	tokenClaim := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := tokenClaim.SignedString([]byte(a.jwtSecret))
	if err != nil {
		return nil, err
	}

	return &response.Login{
		Token:     token,
		ExpiresAt: expireAt.Format(time.RFC3339),
	}, nil
}
