package middleware

import (
	"context"
	"net/http"

	"github.com/zsbahtiar/lionparcel-test/internal/core/model/internalerror"
	"github.com/zsbahtiar/lionparcel-test/internal/pkg/response"
)

type contextKey string

const userIDKey contextKey = "user_id"

func (m *Middleware) AuthUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if len(authHeader) < 1 {
			response.WriteError(w, internalerror.ErrAuthInvalid)
			return
		}
		ctx := r.Context()

		user, err := m.authUsecase.ValidateToken(ctx, authHeader)
		if err != nil {
			response.WriteError(w, err)
			return
		}

		ctx = context.WithValue(ctx, userIDKey, user.ID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (m *Middleware) AuthBackoffice(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if len(authHeader) < 1 {
			response.WriteError(w, internalerror.ErrAuthInvalid)
			return
		}
		ctx := r.Context()

		user, err := m.authUsecase.ValidateToken(ctx, authHeader)
		if err != nil {
			response.WriteError(w, err)
			return
		}

		if !user.IsAdmin {
			response.WriteError(w, internalerror.ErrAuthInvalid)
			return
		}

		ctx = context.WithValue(ctx, userIDKey, user.ID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
