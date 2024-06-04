package services

import "github.com/IntelligenzCodeLab/hacker-news-scraper/data"

type Retriever interface {
	GetItems() ([]data.Item, error)
}
