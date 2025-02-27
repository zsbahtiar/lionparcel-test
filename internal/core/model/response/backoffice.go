package response

type CreateMovie struct {
	ID string `json:"id"`
}

type UpdateMovie struct {
	ID string `json:"id"`
}

type Stats struct {
	MostMovies []struct {
		ID        string `json:"id"`
		Title     string `json:"title"`
		ViewCount int    `json:"view_count"`
	} `json:"most_movies"`
	MostGenres []struct {
		Genre     string `json:"genre"`
		ViewCount int    `json:"view_count"`
	} `json:"most_genres"`

	MostVotedMovies []struct {
		ID        string `json:"id"`
		Title     string `json:"title"`
		VoteCount int    `json:"vote_count"`
	} `json:"most_voted_movies"`
}
