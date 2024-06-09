package services

import (
	"reflect"
	"testing"

	"github.com/IntelligenzCodeLab/hacker-news-scraper/data"
	"github.com/jarcoal/httpmock"
)

const mockWebHtml = `<ol class="stories list "><li id="story_one3oq" data-shortid="one3oq" class="story"><div class="story_liner h-entry"><div class="voters"><a class="upvoter" href="/login"></a><div class="score">21</div></div><div class="details"><span role="heading" aria-level="1" class="link h-cite u-repost-of"><a class="u-url" href="https://www.datagubbe.se/stupidslow/" rel="ugc noreferrer">Stupid Slow: The Perceived Speed of Computers</a></span><span class="tags"><a class="tag tag_performance" title="Performance and optimization" href="/t/performance">performance</a></span><div class="byline"><a href="/~mwcampbell"><img srcset="/avatars/mwcampbell-16.png 1x, /avatars/mwcampbell-32.png 2x" class="avatar" alt="mwcampbell avatar" loading="lazy" decoding="async" src="/avatars/mwcampbell-16.png" width="16" height="16" /></a><span> via </span><a class="u-author h-card" href="/~mwcampbell">mwcampbell</a><span title="2024-06-08 09:04:23 -0500">8 hours ago</span><span> | </span><span class="comments_label"><span> | </span><a role="heading" aria-level="2" href="/s/one3oq/stupid_slow_perceived_speed_computers">              4 comments</a></span></div></div></div></li></ol>`

func TestWebScrapperConnector_GetItems(t *testing.T) {
	type fields struct {
		Url            string
		webContentMock string
		MaxResults     int
	}

	var item1, item2 data.Item
	tests := []struct {
		name          string
		fields        fields
		requestStatus int
		want          []data.Item
		wantErr       bool
	}{
		{name: "One element as response", fields: fields{Url: "http://test1", webContentMock: mockWebHtml, MaxResults: 10}, requestStatus: 200, want: []data.Item{item1}},
		{name: "Two element as response", fields: fields{Url: "http://test2", webContentMock: mockWebHtml, MaxResults: 10}, requestStatus: 200, want: []data.Item{item1, item2}},
		{name: "Error scrapping web page", fields: fields{Url: "http://test3", webContentMock: mockWebHtml, MaxResults: 10}, requestStatus: 200, want: nil, wantErr: true},
		{name: "Error in remote server", fields: fields{Url: "http://test4", webContentMock: mockWebHtml, MaxResults: 10}, requestStatus: 404, want: nil, wantErr: true},
		{name: "Error in network", fields: fields{Url: "http://test5", webContentMock: mockWebHtml, MaxResults: 10}, requestStatus: -1, want: nil, wantErr: true},
	}

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	for _, tt := range tests {
		httpmock.RegisterResponder("GET", tt.fields.Url,
			httpmock.NewStringResponder(200, tt.fields.webContentMock))
		t.Run(tt.name, func(t *testing.T) {
			c := &WebScrapperConnector{Url: tt.fields.Url}
			got, err := c.GetItems(tt.fields.MaxResults)
			if err != nil {
				switch {
				case !tt.wantErr:
					t.Errorf("GetItems() error = %v, wantErr %v", err, tt.wantErr)
					return
				default:
					return
				}
			}
			if reflect.DeepEqual(tt.want, got) {
				t.Errorf("GetItems() want item %v not equals as got %v", tt.want, got)
			}
		})
	}
}
