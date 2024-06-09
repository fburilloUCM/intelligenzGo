package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/IntelligenzCodeLab/hacker-news-scraper/data"
)

const mockWebHtmlBeg = `<!DOCTYPE html>
<html>
<head>
<title>Test Page</title>
</head>
<body>
<ol class="stories list ">`
const mockWebHtmlEnd = `
</ol>
</body>
</html>`
const mockHtmlElemTest1 = `<li id="story_one3oq" data-shortid="one3oq" class="story"><div class="story_liner h-entry"><div class="voters"><a class="upvoter" href="/login"></a><div class="score">21</div></div><div class="details"><span role="heading" aria-level="1" class="link h-cite u-repost-of"><a class="u-url" href="https://www.datagubbe.se/stupidslow/" rel="ugc noreferrer">Stupid Slow: The Perceived Speed of Computers</a></span><span class="tags"><a class="tag tag_performance" title="Performance and optimization" href="/t/performance">performance</a></span><div class="byline"><a href="/~mwcampbell"><img srcset="/avatars/mwcampbell-16.png 1x, /avatars/mwcampbell-32.png 2x" class="avatar" alt="mwcampbell avatar" loading="lazy" decoding="async" src="/avatars/mwcampbell-16.png" width="16" height="16" /></a><span> via </span><a class="u-author h-card" href="/~mwcampbell">mwcampbell</a><span title="2024-06-08 09:04:23 -0500">8 hours ago</span><span> | </span><span class="comments_label"><span> | </span><a role="heading" aria-level="2" href="/s/one3oq/stupid_slow_perceived_speed_computers">              4 comments</a></span></div></div></div></li>`
const mockHtmlElemTest2 = `<li id="story_hilmze" data-shortid="hilmze" class="story"><div class="story_liner h-entry">  <div class="voters">      <a class="upvoter" href="/login"></a>    <div class="score">104</div>  </div>  <div class="details">    <span role="heading" aria-level="1" class="link h-cite u-repost-of">        <a class="u-url" href="https://www.jwz.org/xscreensaver/google.html" rel="ugc noreferrer">XScreenSaver: Google Store Privacy Policy</a>    </span>    <div class="byline">        <span title="2024-06-09 01:22:13 -0500">15 hours ago</span>          <span> | </span>          <span class="comments_label">            <span> | </span>            <a role="heading" aria-level="2" href="/s/hilmze/xscreensaver_google_store_privacy">              10 comments</a>          </span>    </div>  </div></div></li>`

func newUnstartedTestServer() *httptest.Server {
	mux := http.NewServeMux()

	mux.HandleFunc("/test1", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(fmt.Sprintf("%s%s%s", mockWebHtmlBeg, mockHtmlElemTest1, mockWebHtmlEnd)))
	})

	mux.HandleFunc("/test2", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(fmt.Sprintf("%s%s%s%s", mockWebHtmlBeg, mockHtmlElemTest1, mockHtmlElemTest2, mockWebHtmlEnd)))
	})

	mux.HandleFunc("/test3", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(fmt.Sprintf("%s%s%s", mockWebHtmlBeg, `<span>NO ELEMS</span>`, mockWebHtmlEnd)))
	})

	mux.HandleFunc("/test4", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(404)
		w.Write([]byte("Not found\n"))
	})
	return httptest.NewUnstartedServer(mux)
}

func newTestServer() *httptest.Server {
	srv := newUnstartedTestServer()
	srv.Start()
	return srv
}

func TestWebScrapperConnector_GetItems(t *testing.T) {
	ts := newTestServer()
	defer ts.Close()

	type fields struct {
		Url        string
		MaxResults int
	}

	item1Mock := `{"by":"","descendants":4,"id":1,"kids":null,"score":21,"time":0,"title":"Stupid Slow: The Perceived Speed of Computers","type":"","url":""}`
	item2Mock := `{"by":"","descendants":10,"id":2,"kids":null,"score":104,"time":0,"title":"XScreenSaver: Google Store Privacy Policy","type":"","url":""}`
	var item1, item2 data.Item
	getItemFromMock(t, item1Mock, &item1)
	getItemFromMock(t, item2Mock, &item2)
	tests := []struct {
		name          string
		fields        fields
		requestStatus int
		want          []data.Item
		wantErr       bool
	}{
		{name: "One element as response", fields: fields{Url: "/test1", MaxResults: 10}, requestStatus: 200, want: []data.Item{item1}},
		{name: "Two element as response", fields: fields{Url: "/test2", MaxResults: 10}, requestStatus: 200, want: []data.Item{item1, item2}},
		{name: "Error scrapping web page", fields: fields{Url: "/test3", MaxResults: 10}, requestStatus: 200, want: nil, wantErr: true},
		{name: "Error in remote server", fields: fields{Url: "/test4", MaxResults: 10}, requestStatus: 404, want: nil, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &WebScrapperConnector{Url: ts.URL + tt.fields.Url}
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

			maxResults := len(tt.want)
			if !reflect.DeepEqual(tt.want, got[:maxResults]) {
				t.Errorf("GetItems() want item %v not equals as got %v", tt.want, got)
			}
		})
	}
}

func getItemFromMock(t *testing.T, itemMock string, item *data.Item) {
	err := json.Unmarshal([]byte(itemMock), item)
	if err != nil {
		t.Fatalf("Error unmarshaling JSON: %v", err)
	}
}
