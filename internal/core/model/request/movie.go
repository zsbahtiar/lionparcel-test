package request

type GetMovies struct {
	Page   int
	Limit  int
	Search string
}

type VoteMovie struct {
	MovieID string `json:"-" validate:"required"`
	UserID  string `json:"-" validate:"required"`
	Action  string `json:"action" validate:"required,oneof=upvote downvote"`
}

type CreateUserMovieView struct {
	MovieID         string  `json:"-" validate:"required"`
	UserID          string  `json:"-" validate:"required"`
	DurationWatched float64 `json:"duration_watched" validate:"required"`
}
