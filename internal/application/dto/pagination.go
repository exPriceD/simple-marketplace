package dto

type ListingFeedDTO struct {
	Items      []ListingDTO `json:"items"`
	Total      int64        `json:"total"`
	NextOffset int          `json:"next_offset"`
}
