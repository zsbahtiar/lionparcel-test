package handler

import (
	"encoding/json"
	"net/http"

	"github.com/zsbahtiar/lionparcel-test/internal/core/model/request"
	"github.com/zsbahtiar/lionparcel-test/internal/core/module"
)

type authHandler struct {
	authUsecase module.AuthUsecase
}

type AuthHandler interface {
	RegisterUser(w http.ResponseWriter, r *http.Request)
	Login(w http.ResponseWriter, r *http.Request)
	Logout(w http.ResponseWriter, r *http.Request)
}

func NewAuthHandler(authUsecase module.AuthUsecase) AuthHandler {
	return &authHandler{
		authUsecase: authUsecase,
	}
}

func (a *authHandler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	req := &request.Register{}
	defer r.Body.Close()

	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid request body"))
		return
	}

	err := a.authUsecase.RegisterUser(r.Context(), req)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Register Success"))
}

func (a *authHandler) Login(w http.ResponseWriter, r *http.Request) {
	req := &request.Login{}
	defer r.Body.Close()

	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid request body"))
		return
	}

	resp, err := a.authUsecase.Login(r.Context(), req)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

func (a *authHandler) Logout(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")
	if len(token) < 1 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Token is required"))
		return
	}
	req := &request.Logout{Token: token}

	err := a.authUsecase.Logout(r.Context(), req)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Logout Success"))
}
