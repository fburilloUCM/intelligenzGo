package services

import (
	"encoding/json"
	"github.com/IntelligenzCodeLab/hacker-news-scraper/data"
	mock_services "github.com/IntelligenzCodeLab/hacker-news-scraper/services/mock"
	"github.com/golang/mock/gomock"
	"reflect"
	"testing"
)

func TestAggregator_GetItems(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockItemsApiResponse := `[{"by":"testApi1","descendants":12,"id":40540952,"kids":[],"score":746,"time":1717194188,"title":"Test1 item: more that 5","type":"story","url":"https://test/test1"},{"by":"testApi2","descendants":25,"id":40541559,"kids":[],"score":250,"time":1717194188,"title":"1234","type":"story","url":"https://test/test2"},{"by":"testApi3","descendants":55,"id":40540954,"kids":[],"score":111,"time":1717194188,"title":"Test3 item: more that 5","type":"story","url":"https://test/test"}]`
	mockApiFetcher, itemsApiRetriever := createFetcherMock(mockItemsApiResponse, t, ctrl)

	mockItemsWebResponse := `[{"by":"testWeb1","descendants":99,"id":40540952,"kids":[],"score":746,"time":1717194188,"title":"123","type":"story","url":"https://testweb/test1"},{"by":"testWeb2","descendants":1,"id":40541559,"kids":[],"score":250,"time":1717194188,"title":"Test web 2 item: more that 5","type":"story","url":"https://testweb/test2"},{"by":"testWeb3","descendants":80,"id":40540954,"kids":[],"score":1,"time":1717194188,"title":"Test web 3 item: more that 5","type":"story","url":"https://testweb/test3"}]`
	mockWebFetcher, itemsWebRetriever := createFetcherMock(mockItemsWebResponse, t, ctrl)

	// Set up expectations
	mockApiFetcher.EXPECT().GetItems(gomock.Any()).Return(itemsApiRetriever, nil).AnyTimes()
	apiRetriever := SourceConnectors{SourceName: "testApi", Connector: mockApiFetcher}
	mockWebFetcher.EXPECT().GetItems(gomock.Any()).Return(itemsWebRetriever, nil).AnyTimes()
	webRetriever := SourceConnectors{SourceName: "testWeb", Connector: mockWebFetcher}
	wantApiFetchResult := []data.Item{itemsApiRetriever[2], itemsApiRetriever[0], itemsApiRetriever[1]}
	wantWebFetchResult := []data.Item{itemsWebRetriever[2], itemsWebRetriever[1], itemsWebRetriever[0]}
	wantCombinedFetchResult := []data.Item{itemsWebRetriever[2], itemsApiRetriever[2], itemsApiRetriever[0], itemsWebRetriever[1], itemsWebRetriever[0], itemsApiRetriever[1]}

	type fields struct {
		Connectors []SourceConnectors
	}
	type args struct {
		maxItems int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []data.Item
		wantErr bool
	}{
		{name: "Connector API", fields: fields{Connectors: []SourceConnectors{apiRetriever}}, args: args{maxItems: len(wantApiFetchResult)}, want: wantApiFetchResult, wantErr: false},
		{name: "Connector Web", fields: fields{Connectors: []SourceConnectors{webRetriever}}, args: args{maxItems: len(wantWebFetchResult)}, want: wantWebFetchResult, wantErr: false},
		{name: "Combined connectors (API & Web)", fields: fields{Connectors: []SourceConnectors{apiRetriever, webRetriever}}, args: args{maxItems: len(wantApiFetchResult) + len(wantWebFetchResult)}, want: wantCombinedFetchResult, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			agg := &Aggregator{
				Connectors: tt.fields.Connectors,
			}
			got, err := agg.GetItems(tt.args.maxItems)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetItems() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetItems() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func createFetcherMock(mockItemResponse string, t *testing.T, ctrl *gomock.Controller) (*mock_services.MockRetriever, []data.Item) {
	mockApiFetcher := mock_services.NewMockRetriever(ctrl)
	var items []data.Item
	err := json.Unmarshal([]byte(mockItemResponse), &items)
	if err != nil {
		t.Fatalf("Error unmarshaling JSON: %v", err)
	}
	return mockApiFetcher, items
}
