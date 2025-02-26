package handler

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/zsbahtiar/lionparcel-test/internal/core/model/request"
	"github.com/zsbahtiar/lionparcel-test/internal/core/module"
)

type backOfficeHandler struct {
	backofficeUsecase module.BackofficeUsecase
}

type BackofficeHandler interface {
	CreateMovie(w http.ResponseWriter, r *http.Request)
	UpdateMovie(w http.ResponseWriter, r *http.Request)
}

func NewBackofficeHandler(backofficeUsecase module.BackofficeUsecase) BackofficeHandler {
	return &backOfficeHandler{
		backofficeUsecase: backofficeUsecase,
	}
}

func (b *backOfficeHandler) CreateMovie(w http.ResponseWriter, r *http.Request) {
	req := &request.CreateMovie{}

	defer r.Body.Close()

	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		// @Todo: change after all integation
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid request body"))
		return
	}

	resp, err := b.backofficeUsecase.CreateMovie(r.Context(), req)
	if err != nil {
		// @Todo: change after all integation
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Invalid request body"))
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}

func (b *backOfficeHandler) UpdateMovie(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	req := &request.UpdateMovie{
		ID: vars["id"],
	}

	defer r.Body.Close()

	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		// @Todo: change after all integation
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid request body"))
		return
	}

	resp, err := b.backofficeUsecase.UpdateMovice(r.Context(), req)
	if err != nil {
		// @Todo: change after all integation
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Invalid request body"))
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}
