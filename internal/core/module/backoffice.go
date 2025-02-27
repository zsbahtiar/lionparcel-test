package module

import (
	"context"

	"github.com/google/uuid"
	"github.com/zsbahtiar/lionparcel-test/internal/core/model/entity"
	"github.com/zsbahtiar/lionparcel-test/internal/core/model/request"
	"github.com/zsbahtiar/lionparcel-test/internal/core/model/response"
	"github.com/zsbahtiar/lionparcel-test/internal/core/repository"
	"golang.org/x/sync/errgroup"
)

type backOfficeUsecase struct {
	movieRepo repository.MovieRepository
}

type BackofficeUsecase interface {
	CreateMovie(ctx context.Context, req *request.CreateMovie) (*response.CreateMovie, error)
	UpdateMovice(ctx context.Context, req *request.UpdateMovie) (*response.UpdateMovie, error)
	GetStats(ctx context.Context) (*response.Stats, error)
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

func (b *backOfficeUsecase) GetStats(ctx context.Context) (*response.Stats, error) {
	var (
		mostMovies []entity.MovieViewCount
		mostGenres []entity.GenreViewCount
		mostVoted  []entity.MovieVotedCount
	)

	eg, bgCtx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		var err error
		mostMovies, err = b.movieRepo.GetMostViewedMovies(bgCtx, 10)
		return err
	})
	eg.Go(func() error {
		var err error
		mostGenres, err = b.movieRepo.GetMostViewedGenres(bgCtx, 10)
		return err
	})
	eg.Go(func() error {
		var err error
		mostVoted, err = b.movieRepo.GetMostVotedMovies(bgCtx)
		return err
	})
	if err := eg.Wait(); err != nil {
		return nil, err
	}

	statResp := &response.Stats{
		MostMovies: make([]struct {
			ID        string `json:"id"`
			Title     string `json:"title"`
			ViewCount int    `json:"view_count"`
		}, len(mostMovies)),
		MostGenres: make([]struct {
			Genre     string `json:"genre"`
			ViewCount int    `json:"view_count"`
		}, len(mostGenres)),
		MostVotedMovies: make([]struct {
			ID        string `json:"id"`
			Title     string `json:"title"`
			VoteCount int    `json:"vote_count"`
		}, len(mostVoted)),
	}

	for i, m := range mostMovies {
		statResp.MostMovies[i] = struct {
			ID        string `json:"id"`
			Title     string `json:"title"`
			ViewCount int    `json:"view_count"`
		}{
			ID:        m.ID,
			Title:     m.Title,
			ViewCount: m.ViewCount,
		}
	}

	for i, g := range mostGenres {
		statResp.MostGenres[i] = struct {
			Genre     string `json:"genre"`
			ViewCount int    `json:"view_count"`
		}{
			Genre:     g.Genre,
			ViewCount: g.ViewCount,
		}
	}

	for i, v := range mostVoted {
		statResp.MostVotedMovies[i] = struct {
			ID        string `json:"id"`
			Title     string `json:"title"`
			VoteCount int    `json:"vote_count"`
		}{
			ID:        v.ID,
			Title:     v.Title,
			VoteCount: v.VotedCount,
		}
	}
	return statResp, nil
}
