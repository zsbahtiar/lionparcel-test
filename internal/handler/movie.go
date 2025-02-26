package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/zsbahtiar/lionparcel-test/internal/core/model/request"
	"github.com/zsbahtiar/lionparcel-test/internal/core/module"
)

type movieHandler struct {
	movieUsecase module.MovieUsecase
}

type MovieHandler interface {
	GetMovies(w http.ResponseWriter, r *http.Request)
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
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
