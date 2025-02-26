package module

import (
	"context"

	"github.com/google/uuid"
	"github.com/zsbahtiar/lionparcel-test/internal/core/model/entity"
	"github.com/zsbahtiar/lionparcel-test/internal/core/model/request"
	"github.com/zsbahtiar/lionparcel-test/internal/core/model/response"
	"github.com/zsbahtiar/lionparcel-test/internal/core/repository"
)

type backOfficeUsecase struct {
	movieRepo repository.MovieRepository
}

type BackofficeUsecase interface {
	CreateMovie(ctx context.Context, req *request.CreateMovie) (*response.CreateMovie, error)
	UpdateMovice(ctx context.Context, req *request.UpdateMovie) (*response.UpdateMovie, error)
}

func NewBackofficeUsecase(movieRepo repository.MovieRepository) BackofficeUsecase {
	return &backOfficeUsecase{movieRepo: movieRepo}
}

func (b *backOfficeUsecase) CreateMovie(ctx context.Context, req *request.CreateMovie) (*response.CreateMovie, error) {
	movie := &entity.Movie{
		ID:          uuid.New().String(),
		Title:       req.Title,
		Description: req.Description,
		Duration:    req.Duration,
		Artists:     req.Artists,
		Genres:      req.Genres,
		Link:        req.Link,
	}

	err := b.movieRepo.CreateMovie(ctx, movie)
	if err != nil {
		return nil, err
	}

	return &response.CreateMovie{
		ID: movie.ID,
	}, nil
}

func (b *backOfficeUsecase) UpdateMovice(ctx context.Context, req *request.UpdateMovie) (*response.UpdateMovie, error) {
	movie := &entity.Movie{
		ID:          req.ID,
		Title:       req.Title,
		Description: req.Description,
		Duration:    req.Duration,
		Artists:     req.Artists,
		Genres:      req.Genres,
		Link:        req.Link,
	}

	err := b.movieRepo.UpdateMovie(ctx, movie)
	if err != nil {
		return nil, err
	}

	return &response.UpdateMovie{
		ID: movie.ID,
	}, nil
}
