package module

import (
	"context"

	"github.com/zsbahtiar/lionparcel-test/internal/core/model/request"
	"github.com/zsbahtiar/lionparcel-test/internal/core/model/response"
	"github.com/zsbahtiar/lionparcel-test/internal/core/repository"
)

type movieUsecase struct {
	movieRepository repository.MovieRepository
}

type MovieUsecase interface {
	GetMovies(ctx context.Context, req *request.GetMovies) (*response.GetMovies, error)
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
