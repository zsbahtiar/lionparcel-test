package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/zsbahtiar/lionparcel-test/internal/core/model/internalerror"
	"github.com/zsbahtiar/lionparcel-test/internal/core/model/request"
	"github.com/zsbahtiar/lionparcel-test/internal/core/module"
	"github.com/zsbahtiar/lionparcel-test/internal/pkg/middleware"
	"github.com/zsbahtiar/lionparcel-test/internal/pkg/response"
)

type movieHandler struct {
	movieUsecase module.MovieUsecase
}

type MovieHandler interface {
	GetMovies(w http.ResponseWriter, r *http.Request)
	GetMovieView(w http.ResponseWriter, r *http.Request)
	VoteMovie(w http.ResponseWriter, r *http.Request)
	GetVotedMovieOfUser(w http.ResponseWriter, r *http.Request)
	CreateMovieView(w http.ResponseWriter, r *http.Request)
}

func NewMovieHandler(movieUsecase module.MovieUsecase) MovieHandler {
	return &movieHandler{
		movieUsecase: movieUsecase,
	}
}

func (m *movieHandler) GetMovies(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	req := &request.GetMovies{
		Page:   1,
		Limit:  10,
		Search: "",
	}

	if limitStr := query.Get("limit"); limitStr != "" {
		limit, err := strconv.Atoi(limitStr)
		if err == nil && limit > 0 {
			req.Limit = limit
		}
	}

	if pageStr := query.Get("page"); pageStr != "" {
		page, err := strconv.Atoi(pageStr)
		if err == nil && page > 0 {
			req.Page = page
		}
	}

	if search := query.Get("search"); search != "" {
		req.Search = search
	}

	resp, err := m.movieUsecase.GetMovies(r.Context(), req)
	if err != nil {
		response.WriteError(w, err)
		return
	}

	response.WriteSuccess(w, http.StatusOK, resp)
}

func (m *movieHandler) GetMovieView(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	movieID := vars["id"]

	resp, err := m.movieUsecase.GetMovieView(r.Context(), movieID)
	if err != nil {
		response.WriteError(w, err)
		return
	}

	response.WriteSuccess(w, http.StatusOK, resp)
}

func (m *movieHandler) VoteMovie(w http.ResponseWriter, r *http.Request) {
	userID, err := middleware.GetUserID(r.Context())
	if err != nil {
		response.WriteError(w, err)
		return
	}
	mux := mux.Vars(r)
	req := &request.VoteMovie{
		MovieID: mux["id"],
		UserID:  userID,
	}
	fmt.Println(req.MovieID, req.UserID)
	defer r.Body.Close()

	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		response.WriteError(w, internalerror.ErrRequestInvalid)
		return
	}

	err = m.movieUsecase.VoteMovie(r.Context(), req)
	if err != nil {
		response.WriteError(w, err)
		return
	}

	response.WriteSuccess(w, http.StatusOK, nil)
}

func (m *movieHandler) GetVotedMovieOfUser(w http.ResponseWriter, r *http.Request) {
	userID, err := middleware.GetUserID(r.Context())
	if err != nil {
		response.WriteError(w, err)
		return
	}

	resp, err := m.movieUsecase.GetVotedMovieOfUser(r.Context(), userID)
	if err != nil {
		response.WriteError(w, err)
		return
	}

	response.WriteSuccess(w, http.StatusOK, resp)
}

func (m *movieHandler) CreateMovieView(w http.ResponseWriter, r *http.Request) {
	userID, err := middleware.GetUserID(r.Context())
	if err != nil {
		response.WriteError(w, err)
		return
	}
	vars := mux.Vars(r)
	req := &request.CreateUserMovieView{
		UserID:  userID,
		MovieID: vars["id"],
	}

	defer r.Body.Close()

	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		response.WriteError(w, internalerror.ErrRequestInvalid)
		return
	}

	err = m.movieUsecase.CreateUserMovieView(r.Context(), req)
	if err != nil {
		response.WriteError(w, err)
		return
	}

	response.WriteSuccess(w, http.StatusOK, nil)

}
