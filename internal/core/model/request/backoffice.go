package request

/*
example:

	{
		"title": "The Avengers",
		"description": "The Avengers is a team of superheroes appearing in American comic books published by Marvel Comics.",
		"duration": 120,
		"artists": ["Robert Downey Jr.", "Chris Evans", "Mark Ruffalo", "Chris Hemsworth", "Scarlett Johansson", "Jeremy Renner", "Tom Hiddleston", "Clark Gregg", "Cobie Smulders", "Stellan Skarsgård", "Samuel L. Jackson"],
		"genres": ["Action", "Adventure", "Sci-Fi"],
		"link": "https://www.youtube.com/watch?v=eOrNdBpGMv8"
	}
*/
type CreateMovie struct {
	Title       string   `json:"title" validate:"required"`
	Description string   `json:"description" validate:"required"`
	Duration    int      `json:"duration" validate:"required|min=1"`
	Artists     []string `json:"artists" validate:"required"`
	Genres      []string `json:"genres" validate:"required"`
	Link        string   `json:"link" validate:"required|url"`
}

/*
example:

	{
		"title": "The Avengers",
		"description": "The Avengers is a team of superheroes appearing in American comic books published by Marvel Comics.",
		"duration": 120,
		"artists": ["Robert Downey Jr.", "Chris Evans", "Mark Ruffalo", "Chris Hemsworth", "Scarlett Johansson", "Jeremy Renner", "Tom Hiddleston", "Clark Gregg", "Cobie Smulders", "Stellan Skarsgård", "Samuel L. Jackson"],
		"genres": ["Action", "Adventure", "Sci-Fi"],
		"link": "https://www.youtube.com/watch?v=eOrNdBpGMv8"
	}
*/
type UpdateMovie struct {
	ID          string   `json:"id" validate:"required"`
	Title       string   `json:"title" validate:"required"`
	Description string   `json:"description" validate:"required"`
	Duration    int      `json:"duration" validate:"required|min=1"`
	Artists     []string `json:"artists" validate:"required"`
	Genres      []string `json:"genres" validate:"required"`
	Link        string   `json:"link" validate:"required|url"`
}
