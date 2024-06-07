package services

import (
	"cmp"
	"github.com/IntelligenzCodeLab/hacker-news-scraper/data"
	"slices"
)

type SourceConnectors struct {
	SourceName string
	Connector  Retriever
}

type Aggregator struct {
	Connectors []SourceConnectors
}

func (agg *Aggregator) GetItems(maxItems int) ([]data.Item, error) {
	itemsPerSource := maxItems / len(agg.Connectors)
	aggregatedItems := make([]data.Item, 0)
	for _, sourceConnector := range agg.Connectors {
		connector := sourceConnector.Connector
		items, err := connector.GetItems(itemsPerSource)
		if err != nil {
			return nil, err
		}
		if len(items) >= itemsPerSource {
			aggregatedItems = append(aggregatedItems, items[:itemsPerSource]...)
		} else {
			aggregatedItems = append(aggregatedItems, items...)
		}
	}
	aggregatedItems = sortByTitleType(aggregatedItems)
	return aggregatedItems, nil
}

func sortByTitleType(items []data.Item) []data.Item {
	longTitleItems := make([]data.Item, 0)
	shortTitleItems := make([]data.Item, 0)
	for _, item := range items {
		if len(item.Title) < 5 {
			shortTitleItems = append(shortTitleItems, item)
		} else {
			longTitleItems = append(longTitleItems, item)
		}
	}

	slices.SortFunc(longTitleItems, func(a, b data.Item) int {
		return cmp.Compare(a.Descendants, b.Descendants)
	})

	slices.SortFunc(shortTitleItems, func(a, b data.Item) int {
		return cmp.Compare(a.Score, b.Score)
	})
	return append(longTitleItems, shortTitleItems...)
}
