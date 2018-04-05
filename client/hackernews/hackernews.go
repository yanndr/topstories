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
}

// New return a new HackerNews client.
func New() client.Client {
	return &hackernews{}
}

type item struct {
	Title string `json:"title"`
	URL   string `json:"url"`
}

func (hackernews) Get(limit int) (<-chan string, error) {
	b, err := client.Get(topStoriesURL)
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
			defer wg.Done()
			b, err := client.Get(fmt.Sprintf(itemURL, id))
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
		}(id)
	}

	go func() {
		wg.Wait()
		close(resp)
	}()

	return resp, nil

}

func parseIDs(r io.Reader) ([]int, error) {

	body, err := ioutil.ReadAll(r)

	if err != nil {
		return nil, err
	}

	var keys []int
	err = json.Unmarshal(body, &keys)

	if err != nil {
		return nil, err
	}

	return keys, nil
}

func parseItem(r io.Reader) (*item, error) {
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	i := &item{}
	err = json.Unmarshal(b, i)
	if err != nil {
		return nil, err
	}

	return i, nil
}
