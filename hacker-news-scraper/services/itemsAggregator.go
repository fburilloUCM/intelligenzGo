package services

import (
	"cmp"
	"log"
	"slices"
	"strings"
	"sync"

	"github.com/IntelligenzCodeLab/hacker-news-scraper/data"
)

type SourceConnectors struct {
	SourceName string
	Connector  Retriever
}

type Aggregator struct {
	Connectors []SourceConnectors
}

type SourceFetchResult struct {
	Items []data.Item
	Error error
}

func (agg *Aggregator) GetItems(maxItems int) ([]data.Item, error) {
	connectorsNames := make([]string, len(agg.Connectors))
	for i, cnn := range agg.Connectors {
		connectorsNames[i] = cnn.SourceName
	}
	log.Printf("Fetching results from sources: %s", strings.Join(connectorsNames, ","))
	itemsPerSource := maxItems / len(agg.Connectors)
	aggregatedItems := make([]data.Item, 0)
	var wg sync.WaitGroup
	channel := make(chan SourceFetchResult)
	for _, sourceConnector := range agg.Connectors {
		connector := sourceConnector.Connector
		wg.Add(1)
		go func() {
			items, err := connector.GetItems(itemsPerSource)
			channel <- SourceFetchResult{Items: items, Error: err}
			wg.Done()
		}()
	}

	go func() {
		wg.Wait()
		close(channel)
	}()

	for fetchResponse := range channel {
		items := fetchResponse.Items
		err := fetchResponse.Error
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
		return cmp.Compare(b.Descendants, a.Descendants)
	})

	slices.SortFunc(shortTitleItems, func(a, b data.Item) int {
		return cmp.Compare(b.Score, a.Score)
	})
	return append(longTitleItems, shortTitleItems...)
}
