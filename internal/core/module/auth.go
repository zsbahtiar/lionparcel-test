package module

import (
	"context"
	"errors"
	"strings"
	"sync"
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
	userRepository     repository.UserRepository
	jwtSecret          string
	tokenBlacklist     map[string]time.Time
	tokenBlacklistLock sync.RWMutex
}

type AuthUsecase interface {
	RegisterUser(ctx context.Context, req *request.Register) error
	Login(ctx context.Context, req *request.Login) (*response.Login, error)
	Logout(ctx context.Context, req *request.Logout) error
}

func NewAuthUsecase(userRepository repository.UserRepository, jwtSecret string) AuthUsecase {
	usecase := &authUsecase{userRepository: userRepository, jwtSecret: jwtSecret, tokenBlacklist: make(map[string]time.Time)}

	go func() {
		ticker := time.NewTicker(1 * time.Hour)
		defer ticker.Stop()

		for range ticker.C {
			usecase.cleanupBlacklist()
		}
	}()

	return usecase
}

const (
	// @Todo: move to env
	AccessTokenExpiration = 1 * time.Hour
)

var (
	ErrInvalidToken = errors.New("invalid token")
	ErrExpiredToken = errors.New("token has expired")
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

// currently im using in memory, because this is just a test
// in real world, we need to store the token in cache or database
func (a *authUsecase) Logout(ctx context.Context, req *request.Logout) error {
	tokenString := strings.TrimPrefix(req.Token, "Bearer ")

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(a.jwtSecret), nil
	})
	if err != nil {
		return err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return ErrInvalidToken
	}

	exp, ok := claims["exp"].(float64)
	if !ok {
		return ErrInvalidToken
	}

	expiresAt := time.Unix(int64(exp), 0)

	a.tokenBlacklistLock.Lock()
	a.tokenBlacklist[tokenString] = expiresAt
	a.tokenBlacklistLock.Unlock()

	return nil

}

func (a *authUsecase) cleanupBlacklist() {
	now := time.Now()

	a.tokenBlacklistLock.Lock()
	defer a.tokenBlacklistLock.Unlock()

	for token, expiry := range a.tokenBlacklist {
		if now.After(expiry) {
			delete(a.tokenBlacklist, token)
		}
	}
}
