package main

import (
	"fmt"
	"log"

	"github.com/yanndr/topstories/client/hackernews"
)

func main() {
	fmt.Println("topstories dispays the top stories of a news aggregator api. (hakernews)")
	c := hackernews.New()

	resp, err := c.Get(20)

	if err != nil {
		log.Panic(err)
	}

	for r := range resp {
		fmt.Println("resp: ", r)
	}
}
