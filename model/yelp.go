package model

type SearchResponse struct {
	Total      int        `json:"total"`
	Businesses []Business `json:"businesses"`
	Region     Region     `json:"region"`
}

type Business struct {
	Rating       float64     `json:"rating"`
	Price        string      `json:"price"`
	Phone        string      `json:"phone"`
	ID           string      `json:"id"`
	IsClosed     bool        `json:"is_closed"`
	Categories   []Category  `json:"categories"`
	ReviewCount  int         `json:"review_count"`
	Name         string      `json:"name"`
	URL          string      `json:"url"`
	Coordinates  Coordinates `json:"coordinates"`
	ImageURL     string      `json:"image_url"`
	Location     Location    `json:"location"`
	Distance     float32     `json:"distance"`
	Transactions []string    `json:"transactions"`
}

type Category struct {
	Alias string `json:"alias"`
	Title string `json:"title"`
}

type Location struct {
	City     string `json:"city"`
	Country  string `json:"country"`
	Address1 string `json:"address1"`
	Address2 string `json:"address2"`
	Address3 string `json:"address3"`
	State    string `json:"state"`
	ZipCode  string `json:"zip_code"`
}

type Region struct {
	Center Coordinates `json:"center"`
}

type ErrorResponseWrapper struct {
	Error ErrorResponse `json:"error"`
}

type ErrorResponse struct {
	Code        string `json:"code"`
	Description string `json:"description"`
}
