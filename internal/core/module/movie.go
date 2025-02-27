package module

import (
	"context"
	"time"

	"github.com/zsbahtiar/lionparcel-test/internal/core/model/entity"
	"github.com/zsbahtiar/lionparcel-test/internal/core/model/request"
	"github.com/zsbahtiar/lionparcel-test/internal/core/model/response"
	"github.com/zsbahtiar/lionparcel-test/internal/core/repository"
	"golang.org/x/sync/errgroup"
)

type movieUsecase struct {
	movieRepository repository.MovieRepository
}

type MovieUsecase interface {
	GetMovies(ctx context.Context, req *request.GetMovies) (*response.GetMovies, error)
	GetMovieView(ctx context.Context, movieID string) (*response.GetViewMovies, error)
	VoteMovie(ctx context.Context, req *request.VoteMovie) error
	GetVotedMovieOfUser(ctx context.Context, userID string) (*response.GetVotedMovieOfUser, error)
	CreateUserMovieView(ctx context.Context, view *request.CreateUserMovieView) error
}

func NewMovieUsecase(movieRepository repository.MovieRepository) MovieUsecase {
	return &movieUsecase{movieRepository: movieRepository}
}

func (m *movieUsecase) GetMovies(ctx context.Context, req *request.GetMovies) (*response.GetMovies, error) {
	movies, total, err := m.movieRepository.GetMovies(ctx, req)
	if err != nil {
		return nil, err
	}

	getMoviesResponse := &response.GetMovies{
		Total:  total,
		Movies: make([]response.Movie, len(movies)),
	}
	for idx := range movies {
		getMoviesResponse.Movies[idx] = response.Movie{
			ID:          movies[idx].ID,
			Title:       movies[idx].Title,
			Description: movies[idx].Description,
			Duration:    movies[idx].Duration,
			Link:        movies[idx].Link,
			Artists:     movies[idx].Artists,
			Genres:      movies[idx].Genres,
		}
	}

	return getMoviesResponse, nil
}

func (m *movieUsecase) GetMovieView(ctx context.Context, movieID string) (*response.GetViewMovies, error) {
	var (
		movie     *entity.Movie
		totalView int64
	)

	eg, bgCtx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		var err error
		movie, err = m.movieRepository.GetMovie(bgCtx, movieID)
		return err
	})
	eg.Go(func() error {
		var err error
		totalView, err = m.movieRepository.GetViewMovies(bgCtx, movieID)
		return err
	})

	if err := eg.Wait(); err != nil {
		return nil, err
	}

	return &response.GetViewMovies{
		Movie: response.Movie{
			ID:          movie.ID,
			Title:       movie.Title,
			Description: movie.Description,
			Duration:    movie.Duration,
			Link:        movie.Link,
			Artists:     movie.Artists,
			Genres:      movie.Genres,
		},
		TotalViews: totalView,
	}, nil
}

func (m *movieUsecase) VoteMovie(ctx context.Context, req *request.VoteMovie) error {

	if req.Action == "downvote" {
		err := m.movieRepository.DeleteVote(ctx, req.UserID, req.MovieID)
		if err != nil {
			return err
		}
		return nil
	}

	vote := &entity.Vote{
		UserID:  req.UserID,
		MovieID: req.MovieID,
	}

	err := m.movieRepository.CreateVote(ctx, vote)
	if err != nil {
		return err
	}

	return nil
}

func (m *movieUsecase) GetVotedMovieOfUser(ctx context.Context, userID string) (*response.GetVotedMovieOfUser, error) {
	movies, err := m.movieRepository.GetVotedMovieOfUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	getVotedMovieOfUser := &response.GetVotedMovieOfUser{
		Movies: make([]struct {
			ID      string `json:"id"`
			Title   string `json:"title"`
			VotedAt string `json:"voted_at"`
		}, len(movies)),
	}

	for idx := range movies {
		getVotedMovieOfUser.Movies[idx] = struct {
			ID      string `json:"id"`
			Title   string `json:"title"`
			VotedAt string `json:"voted_at"`
		}{
			ID:      movies[idx].ID,
			Title:   movies[idx].Title,
			VotedAt: movies[idx].VotedAt.Format(time.RFC3339),
		}
	}

	return getVotedMovieOfUser, nil
}

func (m *movieUsecase) CreateUserMovieView(ctx context.Context, req *request.CreateUserMovieView) error {
	return m.movieRepository.CreateUserMovieView(ctx, &entity.UserMovieView{
		UserID:          req.UserID,
		MovieID:         req.MovieID,
		DurationWatched: req.DurationWatched,
	})
}
