package data

type ItemId int64

type Item struct {
	By          string `json:"by"`
	Descendants int    `json:"descendants"`
	Id          ItemId `json:"id"`
	Kids        []int  `json:"kids"`
	Score       int    `json:"score"`
	Time        int    `json:"time"`
	Title       string `json:"title"`
	Type        string `json:"type"`
	Url         string `json:"url"`
}
