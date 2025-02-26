package middleware

import "net/http"

func Setup(handler http.Handler) http.Handler {
	handler = corsMiddleware(handler)
	handler = logMiddleware(handler)
	handler = requestIDMiddleware(handler)
	handler = recoverMiddleware(handler)
	return handler
}
