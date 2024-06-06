package services

import (
	"cmp"
	"encoding/json"
	"fmt"
	"github.com/IntelligenzCodeLab/hacker-news-scraper/data"
	"io"
	"log"
	"math"
	"net/http"
	"slices"
	"sync"
)

type APIConnector struct {
	Url              string
	ItemsEndPoint    string
	ItemDataEndPoint string
	MaxResults       int
}

type itemError string

func (s itemError) Error() string { return string(s) }

func (c *APIConnector) GetItems() ([]data.Item, error) {
	reqUrl := fmt.Sprintf("%s/%s.json", c.Url, c.ItemsEndPoint)
	resp, err := http.Get(reqUrl)
	if err != nil {
		log.Printf("Failed to make request: %v", err)
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		log.Printf("Failed to make request: %s", resp.Status)
		var statusError itemError = itemError("Response status: " + resp.Status)
		return nil, statusError
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Failed to read response body: %v", err)
		return nil, err
	}

	var identifiers []data.ItemId
	if err := json.Unmarshal(body, &identifiers); err != nil {
		log.Printf("Failed to unmarshal JSON: %v", err)
	}

	numItems := int(math.Min(float64(len(identifiers)), float64(c.MaxResults)))
	itemChannel := make(chan data.Item)
	waitGroup := sync.WaitGroup{}
	waitGroup.Add(numItems)
	for i := range numItems {
		go c.getItemData(identifiers[i], itemChannel, &waitGroup)
	}
	//wait for all items retrieved
	go func() {
		waitGroup.Wait()
		close(itemChannel)
	}()

	longTitleItems := make([]data.Item, 0)
	shortTitleItems := make([]data.Item, 0)
	for item := range itemChannel {
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

	if len(longTitleItems)+len(shortTitleItems) < numItems {
		var itemsErr itemError = "There has been an error getting some item"
		return longTitleItems, itemsErr
	} else {
		return append(longTitleItems, shortTitleItems...), nil
	}
}

func (c *APIConnector) getItemData(identifier data.ItemId, channel chan data.Item, waitGroup *sync.WaitGroup) {

	reqUrl := fmt.Sprintf("%s/%s/%d.json", c.Url, c.ItemDataEndPoint, identifier)
	resp, err := http.Get(reqUrl)
	if err != nil {
		log.Printf("Failed to make request: %v", err)
		waitGroup.Done()
		return
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Failed to read response body: %v\n", err)
		waitGroup.Done()
		return
	}
	var item data.Item
	if err := json.Unmarshal(body, &item); err != nil {
		log.Printf("Failed to unmarshal JSON: %v\n", err)
		waitGroup.Done()
		return
	}
	channel <- item
	waitGroup.Done()
}
