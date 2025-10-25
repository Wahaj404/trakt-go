package trakt

type Movie struct {
	Title string `json:"title"`
	Year  int    `json:"year"`
	IDs   struct {
		Slug string `json:"slug"`
	} `json:"ids"`
}

type Show struct {
	Title string `json:"title"`
	Year  int    `json:"year"`
	IDs   struct {
		Slug string `json:"slug"`
	} `json:"ids"`
}
