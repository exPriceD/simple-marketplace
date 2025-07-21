package listing

import (
	"time"
)

// Listing — агрегат объявления.
type Listing struct {
	id          int64
	title       Title
	description Description
	imageURL    ImageURL
	price       Price
	authorID    string
	authorLogin string
	createdAt   time.Time
}

// NewListing создаёт новое объявление.
func NewListing(title Title, description Description, img ImageURL, price Price, authorID, authorLogin string, now time.Time) *Listing {
	return &Listing{
		id:          0,
		title:       title,
		description: description,
		imageURL:    img,
		price:       price,
		authorID:    authorID,
		authorLogin: authorLogin,
		createdAt:   now.UTC(),
	}
}

// WithID возвращает копию с установленным ID (используется после INSERT).
func (l *Listing) WithID(id int64, createdAt time.Time) *Listing {
	cp := *l
	cp.id = id
	cp.createdAt = createdAt.UTC()
	return &cp
}

// RehydrateListing восстанавливает объявление из бд.
func RehydrateListing(id int64, title Title, desc Description, img ImageURL, price Price, authorID, authorLogin string, createdAt time.Time) *Listing {
	return &Listing{
		id:          id,
		title:       title,
		description: desc,
		imageURL:    img,
		price:       price,
		authorID:    authorID,
		authorLogin: authorLogin,
		createdAt:   createdAt.UTC(),
	}
}

// Геттеры.
func (l *Listing) ID() int64            { return l.id }
func (l *Listing) Title() Title         { return l.title }
func (l *Listing) Description() string  { return l.description.String() }
func (l *Listing) ImageURL() string     { return l.imageURL.String() }
func (l *Listing) Price() int64         { return l.price.Int64() }
func (l *Listing) AuthorID() string     { return l.authorID }
func (l *Listing) AuthorLogin() string  { return l.authorLogin }
func (l *Listing) CreatedAt() time.Time { return l.createdAt }
