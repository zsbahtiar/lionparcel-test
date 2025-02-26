package handler

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/zsbahtiar/lionparcel-test/internal/core/model/request"
	"github.com/zsbahtiar/lionparcel-test/internal/core/module"
	"github.com/zsbahtiar/lionparcel-test/internal/pkg/response"
)

type movieHandler struct {
	movieUsecase module.MovieUsecase
}

type MovieHandler interface {
	GetMovies(w http.ResponseWriter, r *http.Request)
	GetMovieView(w http.ResponseWriter, r *http.Request)
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
