package module

import (
	"context"

	"github.com/zsbahtiar/lionparcel-test/internal/core/model/request"
	"github.com/zsbahtiar/lionparcel-test/internal/core/model/response"
)

type backOfficeUsecase struct {
}

type BackofficeUsecase interface {
	CreateMovie(ctx context.Context, req *request.CreateMovie) (*response.CreateMovie, error)
	UpdateMovice(ctx context.Context, req *request.UpdateMovie) (*response.UpdateMovie, error)
}

func NewBackofficeUsecase() BackofficeUsecase {
	return &backOfficeUsecase{}
}

func (b *backOfficeUsecase) CreateMovie(ctx context.Context, req *request.CreateMovie) (*response.CreateMovie, error) {
	return &response.CreateMovie{
		ID: "123",
	}, nil
}

func (b *backOfficeUsecase) UpdateMovice(ctx context.Context, req *request.UpdateMovie) (*response.UpdateMovie, error) {
	return &response.UpdateMovie{
		ID: "123",
	}, nil
}
