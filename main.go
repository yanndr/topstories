package main

import (
	"fmt"
	"log"

	"github.com/yanndr/topstories/client/hackernews"
)

func main() {
	c := hackernews.New(20)

	resp, err := c.Get(20)

	if err != nil {
		log.Panic(err)
	}

	for r := range resp {
		if r.Error != nil {
			log.Panic(err)
		}
		fmt.Printf("|%-55s|%-110s|\n", r.Story.Title(), r.Story.URL())
		if r.Error != nil {
			log.Panic(err)
		}
	}
}
