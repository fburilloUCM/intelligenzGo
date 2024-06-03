package main

import (
	"encoding/json"
	"fmt"
	"github.com/IntelligenzCodeLab/hacker-news-scraper/data"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/jarcoal/httpmock"
)

func Test_SendNotification(t *testing.T) {

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	testingIds := []string{"40540952", "40541559"}
	mockIDsResponse := fmt.Sprintf("[%s]", strings.Join(testingIds, ","))
	fmt.Println(mockIDsResponse)
	mockItemResponse1 := `{"by":"nickca","descendants":0,"id":40540952,"kids":[],"score":746,"time":1717194188,"title":"UI elements with a hand-drawn, sketchy look","type":"story","url":"https://wiredjs.com/"}`
	mockItemResponse2 := `{"by":"jer0me","descendants":0,"id":40541559,"kids":[],"score":245,"time":1717200336,"title":"60 kHz (2022)","type":"story","url":"https://ben.page/wwvb"}`
	var testingItem1, testingItem2 data.Item
	json.Unmarshal([]byte(mockItemResponse1), &testingItem1)
	json.Unmarshal([]byte(mockItemResponse1), &testingItem2)
	testingIdsToItems := map[string]data.Item{testingIds[0]: testingItem1, testingIds[1]: testingItem2}

	httpmock.RegisterResponder("GET", strings.Join([]string{hnApiUrl, "topstories.json"}, "/"),
		httpmock.NewStringResponder(200, mockIDsResponse))
	httpmock.RegisterResponder("GET", fmt.Sprintf("%s/item/%s.json", hnApiUrl, testingIds[0]),
		httpmock.NewStringResponder(200, mockItemResponse1))
	httpmock.RegisterResponder("GET", fmt.Sprintf("%s/item/%s.json", hnApiUrl, testingIds[1]),
		httpmock.NewStringResponder(200, mockItemResponse2))

	req, err := http.NewRequest("GET", "/v0/topstories.json", nil)
	if err != nil {
		t.Fatalf("could not create request: %v", err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(SendNotification)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	var items []data.Item
	if err := json.Unmarshal(rr.Body.Bytes(), &items); err != nil {
		t.Fatalf("could not unmarshal response: %v", err)
	}

	if len(items) != 2 {
		t.Fatalf("expected 2 items, got %d", len(items))
	}

	checkedId := testingIdsToItems[testingIds[0]].Id
	if reqItemId, _ := strconv.Atoi(testingIds[0]); checkedId != reqItemId || items[0].Title != "UI elements with a hand-drawn, sketchy look" {
		t.Errorf("unexpected item: %v", items[0])
	}

}
