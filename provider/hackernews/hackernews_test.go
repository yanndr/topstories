package hackernews

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
)

const response string = `{
	"by" : "s3r3nity",
	"descendants" : 802,
	"id" : 16755530,
	"kids" : [ 16761119, 16756858, 16759324, 16755937, 16756474, 16756260, 16756037, 16755857, 16756660, 16756273, 16761037, 16756233, 16758996, 16755931, 16759753, 16756143, 16761093, 16758652, 16755861, 16761190, 16760197, 16758158, 16756533, 16756329, 16760819, 16757399, 16760359, 16757501, 16755910, 16755815, 16760147, 16756516, 16759540, 16760296, 16759655, 16755866, 16756563, 16758602, 16760261, 16755908, 16756679, 16760441, 16756683, 16758810, 16757267, 16761225, 16756094, 16756743, 16757046, 16757519, 16756972, 16757703, 16755954, 16756469, 16757709, 16756266, 16755883, 16756989, 16758406, 16760194, 16757587, 16758952, 16756521, 16758708, 16755918, 16757072, 16756166 ],
	"score" : 958,
	"time" : 1522855070,
	"title" : "test",
	"type" : "story",
	"url" : "https://test.html"
  }`

func handleString(s string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, s)
	}
}

func failAfterOne(s string) func(http.ResponseWriter, *http.Request) {
	i := 0
	m := sync.RWMutex{}
	return func(w http.ResponseWriter, r *http.Request) {
		m.RLock()
		if i > 2 {
			fmt.Fprintf(w, "error")
		}
		m.RUnlock()
		fmt.Fprintf(w, s)
		m.Lock()
		i++
		m.Unlock()
	}
}

func TestNew(t *testing.T) {
	p := New(10)

	hn, ok := p.(*hackernews)

	if !ok {
		t.Fatalf("Expect %T got %T", &hackernews{}, p)
	}

	if hn.sem == nil {
		t.Fatal("Expect non Nil Semaphore")
	}

}

func TestGetStories(t *testing.T) {

	t.Parallel()

	tt := []struct {
		name         string
		n            int
		idsResponse  string
		itemResponse string
		itemFunc     func(string) func(http.ResponseWriter, *http.Request)
		errorIds     bool
		errorItem    bool
	}{
		{name: "happy", n: 5, idsResponse: "[1,2,3,4,5,6,7,8,9,10]", itemResponse: response, itemFunc: handleString, errorIds: false, errorItem: false},
		{name: "badIds", n: 5, idsResponse: "error", itemResponse: response, itemFunc: handleString, errorIds: true, errorItem: false},
		{name: "badItem", n: 5, idsResponse: "[1,2,3,4,5,6,7,8,9,10]", itemResponse: "error", itemFunc: handleString, errorIds: false, errorItem: true},
		{name: "badSecondItem", n: 5, idsResponse: "[1,2,3,4,5,6,7,8,9,10]", itemResponse: response, itemFunc: failAfterOne, errorIds: false, errorItem: true},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			mux := http.NewServeMux()
			srv := httptest.NewServer(mux)

			c := &hackernews{
				sem:           make(chan int, 20),
				topStoriesURL: fmt.Sprintf("%s/%s", srv.URL, "ids"),
				itemURL:       fmt.Sprintf("%s/%s/%s", srv.URL, "items", "%v"),
			}

			mux.HandleFunc("/ids", handleString(tc.idsResponse))
			mux.HandleFunc("/items/", tc.itemFunc(tc.itemResponse))

			resp, err := c.GetStories(tc.n)

			if (err != nil) != tc.errorIds {
				if tc.errorIds {
					t.Fatalf("expected error, got no error")
				}
				t.Fatalf("expected no error, got %v", err)
			}
			if err != nil && tc.errorIds {
				return
			}

			i := 0
			var errs []error
			for r := range resp {
				i++
				if r.Error != nil {
					errs = append(errs, r.Error)
					continue
				}
				if r.Story == nil {
					t.Fatal("error Story not espected to be nil")
				}
				if r.Story.Title() != "test" {
					t.Fatalf("error expected test %v", r.Story.Title())
				}

				if r.Story.URL() != "https://test.html" {
					t.Fatalf("error expected https://test.html got %v", r.Story.URL())
				}
			}

			if i != tc.n {
				t.Fatalf("error, expect %v got %v", tc.n, i)
			}

			if len(errs) > 0 && !tc.errorItem {
				t.Fatalf("expected no error, got %v", errs[0])
			} else if len(errs) == 0 && tc.errorItem {
				t.Fatal("expected  error, got no error")
			}
		})
	}
}
