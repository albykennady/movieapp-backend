package models

type Movie struct {
	ID          int      `json:"id"`
	Title       string   `json:"title"`
	Year        int      `json:"year"`
	Runtime     int      `json:"runtime"`
	MPAARating  string   `json:"mpaa_rating"`
	Description string   `json:"description"`
	Image       string   `json:"image"`
	CreatedAt   string   `json:"-"`
	UpdatedAt   string   `json:"-"`
	Genres      []*Genre `json:"genres,omitempty"`
	GenresArray []int    `json:"genres_array,omitempty"`
}

type MovieUpdate struct {
	Title       *string `json:"title"`
	Year        *int    `json:"year"`
	Runtime     *int    `json:"runtime"`
	MPAARating  *string `json:"mpaa_rating"`
	Description *string `json:"description"`
	Image       *string `json:"image"`
	GenresArray *[]int  `json:"genres_array"`
}

type Genre struct {
	ID        int    `json:"id"`
	Genre     string `json:"genre"`
	Checked   bool   `json:"checked"`
	CreatedAt string `json:"-"`
	UpdatedAt string `json:"-"`
}
