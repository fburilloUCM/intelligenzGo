package main

import (
	"encoding/json"
	"fmt"
	"github.com/IntelligenzCodeLab/hacker-news-scraper/data"
	"github.com/gorilla/mux"
	"io"
	"log"
	"net/http"
)

const hnApiUrl = "https://hacker-news.firebaseio.com/v0"
const maxReturnItems = 10

type ItemId int64

func SendNotification(w http.ResponseWriter, _ *http.Request) {

	reqUrl := fmt.Sprintf("%s/topstories.json", hnApiUrl)
	resp, err := http.Get(reqUrl)
	if err != nil {
		log.Fatalf("Failed to make request: %v", err)
		http.Error(w, "Failed to create request", http.StatusInternalServerError)
		return
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Failed to read response body: %v", err)
	}

	var identifiers []ItemId
	if err := json.Unmarshal(body, &identifiers); err != nil {
		log.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	items := make([]data.Item, 0)
	itemChannel := make(chan data.Item, maxReturnItems)
	for i := range maxReturnItems {
		go addNewItem(identifiers[i], itemChannel)
	}

	//wait for all items retrieved
	checkItem := true
	for checkItem {
		select {
		case item, ok := <-itemChannel:
			if ok {
				items = append(items, item)
			}
		default:
			if len(items) == maxReturnItems {
				checkItem = false
			}
		}
	}

	// Setting the default content-type header to JSON.
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(items)
}

func addNewItem(identifier ItemId, channel chan data.Item) {

	reqUrl := fmt.Sprintf("%s/item/%d.json", hnApiUrl, identifier)
	fmt.Println(identifier)
	resp, err := http.Get(reqUrl)
	if err != nil {
		log.Fatalf("Failed to make request: %v", err)
		return
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Failed to read response body: %v", err)
	}
	var item data.Item
	if err := json.Unmarshal(body, &item); err != nil {
		log.Fatalf("Failed to unmarshal JSON: %v", err)
	}
	channel <- item
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/hacker-news-items", SendNotification).Methods("GET")
	http.Handle("/", r)

	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("could not start server: %v", err)
	}
}
