package client

import (
	"fmt"
	"io"
	"net/http"
)

// Client is an inteface that define a news aggregator method.
type Client interface {
	Get(int) (<-chan Response, error)
}

// Story is an interface that define a story.
type Story interface {
	Title() string
	URL() string
}

// Response represent a response from a story provider.
type Response struct {
	Story Story
	Error error
}

// GetResponseBody execute a Get request and send body response.
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
