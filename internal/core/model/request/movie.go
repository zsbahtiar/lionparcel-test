package request

type GetMovies struct {
	Page   int
	Limit  int
	Search string
}
