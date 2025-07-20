package mapping

import (
	"errors"
	"net/http"
	"strconv"
)

type CreateListingRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	ImageURL    string `json:"image_url"`
	Price       int64  `json:"price"`
}

type ListQuery struct {
	Limit    int
	Offset   int
	SortBy   string
	SortDir  string
	PriceMin *int64
	PriceMax *int64
}

var ErrInvalidQuery = errors.New("invalid_query")

func ParseListQuery(r *http.Request) ListQuery {
	q := r.URL.Query()
	limit, _ := strconv.Atoi(q.Get("limit"))
	offset, _ := strconv.Atoi(q.Get("offset"))
	sortBy := q.Get("sort_by")
	sortDir := q.Get("sort_dir")
	var pmin *int64
	var pmax *int64
	if v := q.Get("price_min"); v != "" {
		if n, err := strconv.ParseInt(v, 10, 64); err == nil {
			pmin = &n
		}
	}
	if v := q.Get("price_max"); v != "" {
		if n, err := strconv.ParseInt(v, 10, 64); err == nil {
			pmax = &n
		}
	}
	return ListQuery{
		Limit:    limit,
		Offset:   offset,
		SortBy:   sortBy,
		SortDir:  sortDir,
		PriceMin: pmin,
		PriceMax: pmax,
	}
}
