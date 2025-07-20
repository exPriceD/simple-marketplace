package listingrepo

import "time"

// ListingRow — структура уровня хранилища.
type ListingRow struct {
	ID          int64     `db:"id"`
	Title       string    `db:"title"`
	Description string    `db:"description"`
	ImageURL    string    `db:"image_url"`
	Price       int64     `db:"price"`
	AuthorID    string    `db:"author_id"`
	CreatedAt   time.Time `db:"created_at"`
}
