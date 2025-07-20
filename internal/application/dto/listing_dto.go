package dto

type ListingDTO struct {
	ID          int64  `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	ImageURL    string `json:"image_url"`
	Price       int64  `json:"price"`
	AuthorID    string `json:"author_id"`
	AuthorLogin string `json:"author_login,omitempty"`
	CreatedAt   string `json:"created_at"`
	IsOwner     bool   `json:"is_owner,omitempty"`
}
