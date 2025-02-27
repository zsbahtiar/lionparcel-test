package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/zsbahtiar/lionparcel-test/internal/core/model/internalerror"
	"github.com/zsbahtiar/lionparcel-test/internal/core/model/request"
	"github.com/zsbahtiar/lionparcel-test/internal/core/module"
	"github.com/zsbahtiar/lionparcel-test/internal/pkg/response"
)

type backOfficeHandler struct {
	backofficeUsecase module.BackofficeUsecase
	validator         *validator.Validate
}

type BackofficeHandler interface {
	CreateMovie(w http.ResponseWriter, r *http.Request)
	UpdateMovie(w http.ResponseWriter, r *http.Request)
	GetStats(w http.ResponseWriter, r *http.Request)
}

func NewBackofficeHandler(backofficeUsecase module.BackofficeUsecase, validator *validator.Validate) BackofficeHandler {
	return &backOfficeHandler{
		backofficeUsecase: backofficeUsecase,
		validator:         validator,
	}
}

func (b *backOfficeHandler) CreateMovie(w http.ResponseWriter, r *http.Request) {
	req := &request.CreateMovie{}

	defer r.Body.Close()

	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		response.WriteError(w, internalerror.ErrRequestInvalid)
		return
	}

	if err := b.validator.Struct(req); err != nil {
		response.WriteError(w, response.New(http.StatusBadRequest, "REQUEST_INVALID", err.Error()))
		return
	}

	resp, err := b.backofficeUsecase.CreateMovie(r.Context(), req)
	if err != nil {
		response.WriteError(w, err)
		return
	}

	response.WriteSuccess(w, http.StatusCreated, resp)
}

func (b *backOfficeHandler) UpdateMovie(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	req := &request.UpdateMovie{
		ID: vars["id"],
	}

	defer r.Body.Close()

	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		response.WriteError(w, internalerror.ErrRequestInvalid)
		return
	}

	if err := b.validator.Struct(req); err != nil {
		response.WriteError(w, response.New(http.StatusBadRequest, "REQUEST_INVALID", err.Error()))
		return
	}
	resp, err := b.backofficeUsecase.UpdateMovice(r.Context(), req)
	if err != nil {
		response.WriteError(w, err)
		return
	}
	response.WriteSuccess(w, http.StatusOK, resp)
}

func (b *backOfficeHandler) GetStats(w http.ResponseWriter, r *http.Request) {
	resp, err := b.backofficeUsecase.GetStats(r.Context())
	if err != nil {
		response.WriteError(w, err)
		return
	}
	response.WriteSuccess(w, http.StatusOK, resp)
}
