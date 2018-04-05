package hackernews

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"

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

func (hackernews) Get(limit int) (<-chan string, error) {
	b, err := getNews(topStoriesURL)
	if err != nil {
		return nil, fmt.Errorf("error on get news ids request: %s", err)
	}
	defer b.Close()
	ids, err := parseIDs(b)
	if err != nil {
		return nil, fmt.Errorf("error parsinf ids response: %s", err)
	}

	resp := make(chan string)

	ids = ids[:limit]

	for k, id := range ids {
		go func(id int) {
			resp <- fmt.Sprint(k)
		}(id)
		log.Println(k)
	}
	return resp, nil

}

func getNews(url string) (io.ReadCloser, error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error on topstories request: %s", err)
	}

	if res.StatusCode >= 400 {
		return nil, fmt.Errorf("err topstories request: %s", res.Status)
	}

	return res.Body, nil

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
