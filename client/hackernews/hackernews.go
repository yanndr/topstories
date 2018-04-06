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
	T string `json:"title"`
	U string `json:"url"`
}

func (i *item) Title() string {
	return i.T
}

func (i *item) URL() string {
	return i.U
}

func (h hackernews) Get(limit int) (<-chan client.Response, error) {
	ids, err := getIDS(h.topStoriesURL)
	if err != nil {
		return nil, fmt.Errorf("error getting ids: %s", err)
	}

	ids = ids[:limit]

	resp := make(chan client.Response)

	wg := sync.WaitGroup{}
	for _, id := range ids {
		wg.Add(1)
		go func(id int) {
			h.sem <- 1
			defer func() {
				log.Println("done")
				wg.Done()
				<-h.sem
			}()
			r := client.Response{}
			item, err := getItem(fmt.Sprintf(h.itemURL, id))
			if err != nil {
				r.Error = err
			}
			r.Story = item
			resp <- r

		}(id)
	}

	go func() {
		wg.Wait()
		close(resp)
	}()

	return resp, nil

}

func getItem(url string) (*item, error) {
	body, err := client.GetResponseBody(url)
	if err != nil {
		return nil, err
	}
	defer body.Close()

	return parseItem(body)
}

func getIDS(url string) ([]int, error) {
	body, err := client.GetResponseBody(url)
	if err != nil {
		return nil, err
	}
	defer body.Close()

	return parseIDs(body)
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
