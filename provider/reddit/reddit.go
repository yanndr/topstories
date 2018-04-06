package reddit

import (
	"github.com/yanndr/topstories/json"
	"github.com/yanndr/topstories/provider"
)

type reddit struct {
	newsURL string
}

type result struct {
	Data data `json:"data"`
}
type data struct {
	Children []children `json:"children"`
}
type children struct {
	Data *childernData `json:"data"`
}

type childernData struct {
	T string `json:"title"`
	U string `json:"url"`
}

func (i *childernData) Title() string {
	return i.T
}

func (i *childernData) URL() string {
	return i.U
}

const (
	newsURL = "https://www.reddit.com/r/golang/new.json?limit=%v"
)

// New returns a new reddit story provider.
func New() provider.StoryProvider {
	return &reddit{
		newsURL: newsURL,
	}
}

func (p *reddit) GetStories(limit int) (<-chan provider.Response, error) {

	var r result

	err := json.UnmarshalFromURL(p.newsURL, &r)
	if err != nil {
		return nil, err
	}

	ch := make(chan provider.Response)

	go func() {
		resp := provider.Response{}
		for _, v := range r.Data.Children {
			resp.Story = v.Data
			ch <- resp
		}
		close(ch)
	}()

	return ch, nil
}
