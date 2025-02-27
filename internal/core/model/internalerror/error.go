package internalerror

import (
	"net/http"

	"github.com/zsbahtiar/lionparcel-test/internal/pkg/response"
)

var (
	ErrMovieNotFound  = response.New(http.StatusNotFound, "MOVIE_NOT_FOUND", "movie not found")
	ErrUserNotFound   = response.New(http.StatusNotFound, "USER_NOT_FOUND", "user not found")
	ErrUserExist      = response.New(http.StatusConflict, "USER_EXIST", "user already exist")
	ErrUnauth         = response.New(http.StatusUnauthorized, "UNAUTHENTICATION", "unauthentication")
	ErrRequestInvalid = response.New(http.StatusBadRequest, "REQUEST_INVALID", "request invalid")
	ErrAuthInvalid    = response.New(http.StatusUnauthorized, "AUTH_INVALID", "authentication invalid")
)
