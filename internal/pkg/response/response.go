package response

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/zsbahtiar/lionparcel-test/internal/pkg/logger"
)

type Response struct {
	Data    any    `json:"data"`
	Code    string `json:"code,omitzero"`
	Message string `json:"message"`
}

type ErrorInfo struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type Error struct {
	StatusCode int    `json:"-"`
	Code       string `json:"code"`
	Message    string `json:"message"`
}

func (e *Error) Error() string {
	return e.Message
}

func New(statusCode int, code, message string) *Error {
	return &Error{
		StatusCode: statusCode,
		Code:       code,
		Message:    message,
	}
}

func WriteSuccess(w http.ResponseWriter, statusCode int, data any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	resp := Response{
		Data:    data,
		Message: "success",
	}

	return json.NewEncoder(w).Encode(resp)
}

func WriteError(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json")

	if e, ok := err.(*Error); ok {
		w.WriteHeader(e.StatusCode)

		resp := Response{
			Code:    e.Code,
			Message: e.Message,
		}

		json.NewEncoder(w).Encode(resp)
		return
	}
	logger.Error(fmt.Sprintf("unknown error: %v", err))

	w.WriteHeader(http.StatusInternalServerError)
	resp := Response{
		Code:    "ERROR",
		Message: "unknown error",
	}
	json.NewEncoder(w).Encode(resp)
}
