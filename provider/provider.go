package provider

import (
	"fmt"
	"io"
	"net/http"
)

// StoryProvider is an inteface that defines a news aggregator method.
type StoryProvider interface {
	GetStories(int) (<-chan Response, error)
}

// Story is an interface that defines a story.
type Story interface {
	Title() string
	URL() string
}

// Response represents a response from a story provider.
type Response struct {
	Story Story
	Error error
}

// GetResponseBody executes a Get request and returns the body response.
func GetResponseBody(url string) (io.ReadCloser, error) {
	res, err := getRequest(url)
	if err != nil {
		return nil, err
	}
	if 200 > res.StatusCode || res.StatusCode > 299 {
		return nil, fmt.Errorf("err with response: %s", res.Status)
	}
	return res.Body, nil
}

func getRequest(url string) (*http.Response, error) {
	c := &http.Client{}
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("error issuing the request: %s", err)
	}
	req.Header.Set("User-Agent", "github.com/yanndr/topStories")

	return c.Do(req)
}
