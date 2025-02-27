package middleware

import (
	"net/http"

	"github.com/zsbahtiar/lionparcel-test/internal/core/module"
)

type Middleware struct {
	authUsecase module.AuthUsecase
}

func New(authUsecase module.AuthUsecase) *Middleware {
	return &Middleware{
		authUsecase: authUsecase,
	}
}

func (m *Middleware) Do(handler http.Handler) http.Handler {
	handler = corsMiddleware(handler)
	handler = logMiddleware(handler)
	handler = requestIDMiddleware(handler)
	handler = recoverMiddleware(handler)
	return handler
}
