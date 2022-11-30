package models

type Category struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type PostCategoryReadCount struct {
	ReadCount  int   `json:"read_count"`
	CategoryID int64 `json:"category_id"`
}
