package entity

import "time"

type Movie struct {
	ID          string   `db:"id"`
	Title       string   `db:"title"`
	Description string   `db:"description"`
	Duration    int      `db:"duration"`
	Artists     []string `db:"artists"`
	Genres      []string `db:"genres"`
	Link        string   `db:"link"`
}

type MovieViewCount struct {
	ID        string `db:"id"`
	Title     string `db:"title"`
	ViewCount int    `db:"view_count"`
}

type GenreViewCount struct {
	Genre     string `db:"genre"`
	ViewCount int    `db:"view_count"`
}

type MovieVotedCount struct {
	ID         string `db:"id"`
	Title      string `db:"title"`
	VotedCount int    `db:"voted_count"`
}

type Vote struct {
	UserID  string `db:"user_id"`
	MovieID string `db:"movie_id"`
}

type UserMovieVote struct {
	ID      string    `db:"id"`
	Title   string    `db:"title"`
	VotedAt time.Time `db:"voted_at"`
}
