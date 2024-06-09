package services

import "github.com/IntelligenzCodeLab/hacker-news-scraper/data"

//go:generate mockgen -source=itemRetriever.go -destination=mock/itemRetriever.go

type Retriever interface {
	GetItems(maxItems int) ([]data.Item, error)
}
