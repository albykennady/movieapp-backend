package models

type Movie struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Year        int    `json:"year"`
	Runtime     int    `json:"runtime"`
	MPAARating  string `json:"mpaa_rating"`
	Description string `json:"description"`
	Image       string `json:"image"`
	CreatedAt   string `json:"-"`
	UpdatedAt   string `json:"-"`
}
