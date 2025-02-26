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
}

func NewBackofficeUsecase() BackofficeUsecase {
	return &backOfficeUsecase{}
}

func (b *backOfficeUsecase) CreateMovie(ctx context.Context, req *request.CreateMovie) (*response.CreateMovie, error) {
	return &response.CreateMovie{
		ID: "123",
	}, nil
}
