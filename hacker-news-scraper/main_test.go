package main

import (
	"encoding/json"
	"github.com/IntelligenzCodeLab/hacker-news-scraper/data"
	"github.com/IntelligenzCodeLab/hacker-news-scraper/services"
	mock_services "github.com/IntelligenzCodeLab/hacker-news-scraper/services/mock"
	"github.com/golang/mock/gomock"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRetrieveHackerNewsItems(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockFetcher := mock_services.NewMockRetriever(ctrl)
	mockItemResponse := `[{"by":"nickca","descendants":0,"id":40540952,"kids":[],"score":746,"time":1717194188,"title":"UI elements with a hand-drawn, sketchy look","type":"story","url":"https://wiredjs.com/"},{"by":"jer0me","descendants":0,"id":40541559,"kids":[],"score":245,"time":1717200336,"title":"60 kHz (2022)","type":"story","url":"https://ben.page/wwvb"}]`
	var items []data.Item
	err := json.Unmarshal([]byte(mockItemResponse), &items)
	if err != nil {
		t.Fatalf("Error unmarshaling JSON: %v", err)
	}

	// Set up expectations
	mockFetcher.EXPECT().GetItems(maxReturnItems).Return(items, nil)

	// Create the handler with the mock fetcher
	handler := BuildItemsRetrieverHandler(map[string]services.Retriever{"test": mockFetcher})

	req, err := http.NewRequest("GET", "/ids", nil)
	if err != nil {
		t.Fatalf("could not create request: %v", err)
	}

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	var scrapedResult []data.ScraperResponse
	if err := json.Unmarshal(rr.Body.Bytes(), &scrapedResult); err != nil {
		t.Fatalf("could not unmarshal response: %v", err)
	}

	if len(scrapedResult) != 2 {
		t.Fatalf("expected 2 items, got %d", len(scrapedResult))
	}

	oneExpectedTitle := items[0].Title
	anotherExpectedTitle := items[0].Title
	wantItemFound := false
	for _, wantItem := range []string{oneExpectedTitle, anotherExpectedTitle} {
		for _, gotItem := range scrapedResult {
			if gotItem.Title == wantItem {
				wantItemFound = true
			}
		}
		if !wantItemFound {
			t.Errorf("GetItems() want item %v not present in %v", scrapedResult, wantItem)
		}
		wantItemFound = false
	}

}
