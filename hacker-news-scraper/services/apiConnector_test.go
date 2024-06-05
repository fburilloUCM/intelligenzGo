package services

import (
	"encoding/json"
	"fmt"
	"github.com/IntelligenzCodeLab/hacker-news-scraper/data"
	"github.com/jarcoal/httpmock"
	"reflect"
	"strings"
	"testing"
)

func TestAPIConnector_GetItems(t *testing.T) {
	type fields struct {
		Url              string
		ItemsEndPoint    string
		ItemDataEndPoint string
		MaxResults       int
	}
	mockItemResponse1 := `{"by":"nickca","descendants":0,"id":40540952,"kids":[],"score":746,"time":1717194188,"title":"UI elements with a hand-drawn, sketchy look","type":"story","url":"https://wiredjs.com/"}`
	mockItemResponse2 := `{"by":"jer0me","descendants":0,"id":40541559,"kids":[],"score":245,"time":1717200336,"title":"60 kHz (2022)","type":"story","url":"https://ben.page/wwvb"}`
	var item1, item2 data.Item
	err := json.Unmarshal([]byte(mockItemResponse1), &item1)
	if err != nil {
		t.Fatalf("Error unmarshaling JSON: %v", err)
	}
	err2 := json.Unmarshal([]byte(mockItemResponse1), &item2)
	if err2 != nil {
		t.Fatalf("Error unmarshaling JSON: %v", err2)
	}
	testingIds := []string{"40540952", "40541559"}
	testingIdsToItems := map[string]string{testingIds[0]: mockItemResponse1, testingIds[1]: mockItemResponse2}
	tests := []struct {
		name          string
		fields        fields
		requestStatus int
		want          []data.Item
		wantErr       bool
	}{
		{name: "One element as response", fields: fields{Url: "http://test", ItemsEndPoint: "ids", ItemDataEndPoint: "item", MaxResults: 10}, requestStatus: 200, want: []data.Item{item1}},
		{name: "Two element as response", fields: fields{Url: "http://test", ItemsEndPoint: "ids", ItemDataEndPoint: "item", MaxResults: 10}, requestStatus: 200, want: []data.Item{item1, item2}},
		{name: "Error marshaling an item data response", fields: fields{Url: "http://test", ItemsEndPoint: "ids", ItemDataEndPoint: "error", MaxResults: 10}, requestStatus: 200, want: nil, wantErr: true},
		{name: "Error in remote server", fields: fields{Url: "http://test", ItemsEndPoint: "ids", ItemDataEndPoint: "error", MaxResults: 10}, requestStatus: 404, want: nil, wantErr: true},
		{name: "Error in network", fields: fields{Url: "http://test", ItemsEndPoint: "ids", ItemDataEndPoint: "error", MaxResults: 10}, requestStatus: -1, want: nil, wantErr: true},
	}

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	for i, tt := range tests {
		var mockIDsResponse string
		switch {
		case i == 0:
			mockIDsResponse = fmt.Sprintf("[%s]", testingIds[i])
		case i > 0 && tt.requestStatus == 200:
			mockIDsResponse = fmt.Sprintf("[%s]", strings.Join(testingIds, ","))
		default:
			mockIDsResponse = "Response invalid"
		}
		if responseStatus := tt.requestStatus; responseStatus > 0 {
			httpmock.RegisterResponder("GET", strings.Join([]string{tt.fields.Url, tt.fields.ItemsEndPoint + ".json"}, "/"),
				httpmock.NewStringResponder(responseStatus, mockIDsResponse))
		} else {
			httpmock.RegisterResponder("GET", strings.Join([]string{tt.fields.Url, tt.fields.ItemsEndPoint + ".json"}, "/"),
				httpmock.NewErrorResponder(fmt.Errorf("simulated network error")))
		}

		if i < 2 {
			for j := range i + 1 {
				testingIdItem := testingIds[j]
				httpmock.RegisterResponder("GET", fmt.Sprintf("%s/item/%s.json", tt.fields.Url, testingIdItem),
					httpmock.NewStringResponder(200, testingIdsToItems[testingIdItem]))
			}
		} else {
			for j := range len(testingIds) {
				testingIdItem := testingIds[j]
				httpmock.RegisterResponder("GET", fmt.Sprintf("%s/error/%s.json", tt.fields.Url, testingIdItem),
					httpmock.NewStringResponder(200, "asdf"))
			}
		}
		t.Run(tt.name, func(t *testing.T) {
			c := &APIConnector{
				Url:              tt.fields.Url,
				ItemsEndPoint:    tt.fields.ItemsEndPoint,
				ItemDataEndPoint: tt.fields.ItemDataEndPoint,
				MaxResults:       tt.fields.MaxResults,
			}
			got, err := c.GetItems()
			if err != nil {
				switch {
				case !tt.wantErr:
					t.Errorf("GetItems() error = %v, wantErr %v", err, tt.wantErr)
					return
				default:
					return
				}
			}
			wantItemFound := false
			for _, wantItem := range tt.want {
				for _, gotItem := range got {
					if reflect.DeepEqual(wantItem, gotItem) {
						wantItemFound = true
					}
				}
				if !wantItemFound {
					t.Errorf("GetItems() want item %v not present in %v", got, wantItem)
				}
				wantItemFound = false
			}
		})
	}
}
