package entity

type Movie struct {
	ID          string   `db:"id"`
	Title       string   `db:"title"`
	Description string   `db:"description"`
	Duration    int      `db:"duration"`
	Artists     []string `db:"artists"`
	Genres      []string `db:"genres"`
	Link        string   `db:"link"`
}
