package models

type PostType struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Quota int    `json:"quota"`
}
