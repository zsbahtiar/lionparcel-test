package middleware

import (
	"net/http"

	"github.com/google/uuid"
)

func requestIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := uuid.New().String()
		w.Header().Set("X-Request-ID", requestID)
		ctx := r.Context()
		// You could also store the requestID in the context if needed
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
