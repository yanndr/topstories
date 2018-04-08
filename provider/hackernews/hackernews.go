package hackernews

import (
	"fmt"
	"sync"

	"github.com/yanndr/topstories/json"
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
	var ids []int
	err := json.UnmarshalFromURL(h.topStoriesURL, &ids)
	if err != nil {
		return nil, fmt.Errorf("error getting ids: %s", err)
	}

	if limit < len(ids) {
		ids = ids[:limit]
	}

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
			item := &item{}
			err := json.UnmarshalFromURL(fmt.Sprintf(h.itemURL, id), item)
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
