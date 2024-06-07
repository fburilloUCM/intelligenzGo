package services

import (
	"fmt"
	"github.com/IntelligenzCodeLab/hacker-news-scraper/data"
	"github.com/gocolly/colly"
	"strconv"
)

type WebScrapperConnector struct {
	Url string
}

func (ws *WebScrapperConnector) GetItems(maxItems int) ([]data.Item, error) {
	items := make([]data.Item, maxItems)
	collector := colly.NewCollector()
	collector.OnHTML("ol.stories.list", func(elemList *colly.HTMLElement) {
		elemList.ForEach("li.story", func(i int, specListElement *colly.HTMLElement) {
			if i >= maxItems {
				return
			}
			title := specListElement.ChildText("div.h-entry .details .link a")
			scoreText := specListElement.ChildText("div.h-entry .voters .score")
			commentsText := specListElement.ChildText("div.h-entry .details .byline .comments_label a")
			score, _ := strconv.Atoi(scoreText)
			comments, _ := extractComments(commentsText)
			items[i] = data.Item{Id: data.ItemId(i + 1), Title: title, Descendants: comments, Score: score}
		})
	})
	err := collector.Visit(ws.Url)
	if err != nil {
		return nil, err
	}
	return items, nil
}

func extractComments(s string) (int, error) {
	var number int
	n, err := fmt.Sscanf(s, "%d comments", &number)
	if err != nil {
		return 0, fmt.Errorf("failed to scan number: %v", err)
	}
	if n != 1 {
		return 0, fmt.Errorf("no match found")
	}
	return number, nil
}
