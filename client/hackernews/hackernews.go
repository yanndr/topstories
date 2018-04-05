package hackernews

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"sync"

	"github.com/yanndr/topstories/client"
)

const (
	topStoriesURL = "https://hacker-news.firebaseio.com/v0/topstories.json"
	itemURL       = "https://hacker-news.firebaseio.com/v0/item/%v.json"
)

type hackernews struct {
	sem                    chan int
	topStoriesURL, itemURL string
}

// New return a new HackerNews client.
func New(maxGoRoutine int) client.Client {
	return &hackernews{
		sem:           make(chan int, maxGoRoutine),
		topStoriesURL: topStoriesURL,
		itemURL:       itemURL,
	}
}

type item struct {
	Title string `json:"title"`
	URL   string `json:"url"`
}

func (h hackernews) Get(limit int) (<-chan string, error) {
	b, err := client.Get(h.topStoriesURL)
	if err != nil {
		return nil, fmt.Errorf("error on get news ids request: %s", err)
	}
	defer b.Close()
	ids, err := parseIDs(b)
	if err != nil {
		return nil, fmt.Errorf("error parsing ids response: %s", err)
	}

	resp := make(chan string)

	ids = ids[:limit]
	wg := sync.WaitGroup{}
	for _, id := range ids {
		wg.Add(1)
		go func(id int) {
			h.sem <- 1
			defer wg.Done()
			b, err := client.Get(fmt.Sprintf(h.itemURL, id))
			if err != nil {
				log.Printf("need to handle this error: %s", err)
				return
			}
			item, err := parseItem(b)
			if err != nil {
				log.Printf("need to handle this error: %s", err)
				return
			}
			resp <- item.Title
			<-h.sem
		}(id)
	}

	go func() {
		wg.Wait()
		close(resp)
	}()

	return resp, nil

}

func parseIDs(r io.Reader) ([]int, error) {
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("parsing ids failed readall error: %v", err)
	}

	var keys []int
	if len(b) == 0 {
		return keys, nil
	}

	err = json.Unmarshal(b, &keys)
	if err != nil {
		return nil, fmt.Errorf("parsing ids failed Unmasrshal error: %v", err)
	}

	return keys, nil
}

func parseItem(r io.Reader) (*item, error) {
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("parsing item failed readall error: %v", err)
	}

	i := &item{}
	if len(b) == 0 {
		return i, nil
	}

	err = json.Unmarshal(b, i)
	if err != nil {
		return nil, fmt.Errorf("parsing item failed Unmasrshal error: %v", err)
	}

	return i, nil
}
