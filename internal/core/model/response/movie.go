package response

type Movie struct {
	ID          string   `json:"id"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Duration    int      `json:"duration"`
	Artists     []string `json:"artists"`
	Genres      []string `json:"genres"`
	Link        string   `json:"link"`
}

type GetMovies struct {
	Movies []Movie `json:"movies"`
	Total  int64   `json:"total"`
}
