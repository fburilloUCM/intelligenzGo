.PHONY: all test build run

all: test build

test:
	go test ./...

build:
	go build -o intelligenzGo

build-container:
	docker build -t hacker-news-scrapper .

run:
	./intelligenzGo

run-container:
	docker run -p 8080:8080 hacker-news-scrapper