package middleware

import (
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/zsbahtiar/lionparcel-test/internal/pkg/logger"
)

func recoverMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				logger.Error(fmt.Sprintf("Panic recovered: %v\nStack trace: %s", err, debug.Stack()))
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(`{"error": "Internal server error"}`))
			}
		}()
		next.ServeHTTP(w, r)
	})
}
