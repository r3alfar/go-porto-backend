package models

type Movie struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	ReleaseDate int    `json:"release_date"`
	RunTime     int    `json:"runtime"`
	MPAARating  string `json:"mpaa_rating"`
	Description string `json:"description"`
	Image       string `json:"image"`
	CreatedAt   int    `json:"created_at"`
	UpdatedAt   int    `json:"updated_at"`
}
