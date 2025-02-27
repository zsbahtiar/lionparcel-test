package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/zsbahtiar/lionparcel-test/internal/core/model/internalerror"
	"github.com/zsbahtiar/lionparcel-test/internal/core/model/request"
	"github.com/zsbahtiar/lionparcel-test/internal/core/module"
	"github.com/zsbahtiar/lionparcel-test/internal/pkg/response"
)

type authHandler struct {
	authUsecase module.AuthUsecase
	validator   *validator.Validate
}

type AuthHandler interface {
	RegisterUser(w http.ResponseWriter, r *http.Request)
	Login(w http.ResponseWriter, r *http.Request)
	Logout(w http.ResponseWriter, r *http.Request)
}

func NewAuthHandler(authUsecase module.AuthUsecase, validator *validator.Validate) AuthHandler {
	return &authHandler{
		authUsecase: authUsecase,
		validator:   validator,
	}
}

func (a *authHandler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	req := &request.Register{}
	defer r.Body.Close()

	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		response.WriteError(w, internalerror.ErrRequestInvalid)
		return
	}

	if err := a.validator.Struct(req); err != nil {
		response.WriteError(w, response.New(http.StatusBadRequest, "REQUEST_INVALID", err.Error()))
		return
	}

	err := a.authUsecase.RegisterUser(r.Context(), req)
	if err != nil {
		response.WriteError(w, err)
		return
	}

	response.WriteSuccess(w, http.StatusCreated, nil)
}

func (a *authHandler) Login(w http.ResponseWriter, r *http.Request) {
	req := &request.Login{}
	defer r.Body.Close()

	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		response.WriteError(w, internalerror.ErrRequestInvalid)
		return
	}

	if err := a.validator.Struct(req); err != nil {
		response.WriteError(w, response.New(http.StatusBadRequest, "REQUEST_INVALID", err.Error()))
		return
	}

	resp, err := a.authUsecase.Login(r.Context(), req)
	if err != nil {
		response.WriteError(w, err)
		return
	}

	response.WriteSuccess(w, http.StatusOK, resp)
}

func (a *authHandler) Logout(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")
	req := &request.Logout{Token: token}

	if err := a.validator.Struct(req); err != nil {
		response.WriteError(w, response.New(http.StatusBadRequest, "REQUEST_INVALID", err.Error()))
		return
	}

	err := a.authUsecase.Logout(r.Context(), req)
	if err != nil {
		response.WriteError(w, err)
		return
	}
	response.WriteSuccess(w, http.StatusOK, nil)
}
