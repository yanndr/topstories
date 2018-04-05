package hackernews

import "github.com/yanndr/topstories/client"

type hackernews struct {
}

// New return a new HackerNews client.
func New() client.Client {
	return &hackernews{}
}

func (hackernews) Get() {
}
