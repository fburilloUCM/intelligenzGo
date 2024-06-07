package main

import (
	"encoding/json"
	"fmt"
	"github.com/IntelligenzCodeLab/hacker-news-scraper/data"
	"github.com/IntelligenzCodeLab/hacker-news-scraper/services"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
)

const hackerNewsName = "Hacker News"
const hnApiUrl = "https://hacker-news.firebaseio.com/v0"
const lobstersName = "Lobsters"
const lobstersWebUrl = "https://lobste.rs/"
const maxReturnItems = 30

func BuildItemsRetrieverHandler(sourcesRetrievers map[string]services.Retriever) http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		connectors := make([]services.SourceConnectors, 0)
		for sourceName, retriever := range sourcesRetrievers {
			connectors = append(connectors, services.SourceConnectors{SourceName: sourceName, Connector: retriever})
		}
		aggregator := services.Aggregator{Connectors: connectors}
		items, err := aggregator.GetItems(maxReturnItems)
		if err != nil {
			log.Printf("Failed to get Connector: %v\n", err)
			http.Error(w, "Error obtaining required data", http.StatusInternalServerError)
			return
		}

		// Setting the default content-type header to JSON.
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		response := make([]data.ScraperResponse, len(items))
		for i, item := range items {
			title := item.Title
			num := i + 1
			comments := item.Descendants
			id := int(item.Id)
			score := item.Score
			fmt.Printf("%d - Title(%d): %s (%d comments, score %d). Id: %d\n", num, len(title), title, comments, score, id)
			response[i] = data.ScraperResponse{Order: num, Id: strconv.Itoa(id), Title: title, Url: item.Url, Comments: comments, Score: score}
		}
		err = json.NewEncoder(w).Encode(response)
		if err != nil {
			log.Printf("Failed to build response: %v", err)
			http.Error(w, "Error building service response", http.StatusInternalServerError)
			return
		}
	}
}

func main() {
	r := mux.NewRouter()

	var hackerNewsRetriever services.Retriever
	hackerNewsRetriever = &services.APIConnector{Url: hnApiUrl, ItemsEndPoint: "topstories", ItemDataEndPoint: "item"}
	hackerRankRetrieveHandler := BuildItemsRetrieverHandler(map[string]services.Retriever{hackerNewsName: hackerNewsRetriever})
	r.HandleFunc("/hacker-news-items", hackerRankRetrieveHandler).Methods("GET")

	lobstersRetriever := &services.WebScrapperConnector{Url: lobstersWebUrl}
	lobstersWebScrapeRetrieveHandler := BuildItemsRetrieverHandler(map[string]services.Retriever{lobstersName: lobstersRetriever})
	r.HandleFunc("/lobsters-items", lobstersWebScrapeRetrieveHandler).Methods("GET")

	combinedRetrieverHandler := BuildItemsRetrieverHandler(map[string]services.Retriever{hackerNewsName: hackerNewsRetriever, lobstersName: lobstersRetriever})
	r.HandleFunc("/combine-sources-items", combinedRetrieverHandler).Methods("GET")
	http.Handle("/", r)

	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("could not start server: %v", err)
	}
}
