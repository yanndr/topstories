package client

import (
	"fmt"
	"io"
	"net/http"
)

// Client is an inteface that define a news aggregator method.
type Client interface {
	Get(int) (<-chan string, error)
}

// Get execute a Get request and send the body.
func Get(url string) (io.ReadCloser, error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error issuing the request: %s", err)
	}

	if res.StatusCode >= 400 {
		return nil, fmt.Errorf("err with status: %s", res.Status)
	}

	return res.Body, nil
}
