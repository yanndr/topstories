package reddit

import (
	"github.com/yanndr/topstories/json"
	"github.com/yanndr/topstories/provider"
)

type reddit struct {
	newsURL string
}

type childrenData struct {
	T string `json:"title"`
	U string `json:"url"`
}

func (i *childrenData) Title() string {
	return i.T
}

func (i *childrenData) URL() string {
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

	var r struct {
		Data struct {
			Children []struct {
				Data *childrenData `json:"data"`
			} `json:"children"`
		} `json:"data"`
	}

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
