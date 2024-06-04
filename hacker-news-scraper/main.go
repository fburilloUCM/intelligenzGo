package main

import (
	"encoding/json"
	"github.com/IntelligenzCodeLab/hacker-news-scraper/data"
	"github.com/IntelligenzCodeLab/hacker-news-scraper/services"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

const hnApiUrl = "https://hacker-news.firebaseio.com/v0"
const maxReturnItems = 10

func RetrieveHackerNewsItems(w http.ResponseWriter, _ *http.Request) {

	var retriever services.Retriever
	retriever = &services.APIConnector{Url: hnApiUrl, MaxResults: maxReturnItems, ItemsEndPoint: "topstories", ItemDataEndPoint: "item"}
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
		response[i] = data.ScraperResponse{Title: item.Title}
	}
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Printf("Failed to build response: %v", err)
		http.Error(w, "Error building service response", http.StatusInternalServerError)
		return
	}
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/hacker-news-items", RetrieveHackerNewsItems).Methods("GET")
	http.Handle("/", r)

	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("could not start server: %v", err)
	}
}
