package hackernews

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"sync"

	"github.com/yanndr/topstories/provider"
)

const (
	topStoriesURL = "https://hacker-news.firebaseio.com/v0/topstories.json"
	itemURL       = "https://hacker-news.firebaseio.com/v0/item/%v.json"
)

type hackernews struct {
	sem                    chan int
	topStoriesURL, itemURL string
}

// New returns a new HackerNews story provider.
func New(maxGoRoutine int) provider.StoryProvider {
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

func (h hackernews) GetStories(limit int) (<-chan provider.Response, error) {
	ids, err := getIDS(h.topStoriesURL)
	if err != nil {
		return nil, fmt.Errorf("error getting ids: %s", err)
	}

	ids = ids[:limit]

	resp := make(chan provider.Response)

	wg := sync.WaitGroup{}
	for _, id := range ids {
		wg.Add(1)
		go func(id int) {
			h.sem <- 1

			defer func() {
				wg.Done()
				<-h.sem
			}()

			r := provider.Response{}
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
	body, err := provider.GetResponseBody(url)
	if err != nil {
		return nil, err
	}
	defer body.Close()

	i := &item{}
	err = parse(body, &i)
	if err != nil {
		return nil, err
	}
	return i, err
}

func getIDS(url string) ([]int, error) {
	body, err := provider.GetResponseBody(url)
	if err != nil {
		return nil, err
	}
	defer body.Close()

	var keys []int
	err = parse(body, &keys)
	if err != nil {
		return nil, err
	}

	return keys, err
}

func parse(r io.Reader, o interface{}) error {
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return fmt.Errorf("parsing failed readall error: %v", err)
	}

	if len(b) == 0 {
		return nil
	}

	err = json.Unmarshal(b, o)
	if err != nil {
		return fmt.Errorf("parsing failed Unmasrshal error: %v", err)
	}

	return nil
}
