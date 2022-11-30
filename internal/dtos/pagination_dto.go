package dtos

type PaginationResponse struct {
	PerPage     int   `json:"per_page"`
	CurrentPage int   `json:"current_page"`
	TotalRows   int64 `json:"total"`
	TotalPages  int64 `json:"total_pages"`
}
