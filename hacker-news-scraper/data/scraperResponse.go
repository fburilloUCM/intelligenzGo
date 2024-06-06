package data

type ScraperResponse struct {
	Order    int    `json:"order"`
	Id       string `json:"id"`
	Title    string `json:"title"`
	Url      string `json:"url"`
	Comments int    `json:"comments"`
	Score    int    `json:"score"`
}
