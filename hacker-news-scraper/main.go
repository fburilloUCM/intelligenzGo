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

const hnApiUrl = "https://hacker-news.firebaseio.com/v0"
const maxReturnItems = 30

func RetrieveHackerNewsItems(retriever services.Retriever) http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {

		items, err := retriever.GetItems()
		if err != nil {
			log.Printf("Failed to get Items: %v\n", err)
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

	var retriever services.Retriever
	retriever = &services.APIConnector{Url: hnApiUrl, MaxResults: maxReturnItems, ItemsEndPoint: "topstories", ItemDataEndPoint: "item"}
	hackerRankRetrieveHandler := RetrieveHackerNewsItems(retriever)
	r.HandleFunc("/hacker-news-items", hackerRankRetrieveHandler).Methods("GET")
	http.Handle("/", r)

	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("could not start server: %v", err)
	}
}
