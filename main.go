// MIT License

// Copyright (c) 2018 Yann Druffin

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

// Console application to display stories form news aggregators
package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/yanndr/topstories/provider"
	"github.com/yanndr/topstories/provider/hackernews"
	"github.com/yanndr/topstories/provider/reddit"
)

func main() {
	var (
		output *os.File
		sw     provider.StoryWriter
		err    error
	)

	csvPtr := flag.Bool("csv", false, "Save the result to a csv file.")
	prPtr := flag.String("p", "hackernews", "Stories provider: hackernews or reddit")
	path := flag.String("o", "outupt.csv", "output file name")
	n := flag.Int("n", 20, "number of stories to display")
	c := flag.Uint("c", 20, "max concurency allowed")
	flag.Parse()

	if *c < 1 {
		fmt.Println("Cannot have less than 1 go routine allowed.")
		os.Exit(1)
	}

	if *csvPtr {
		output, err = os.Create(*path)
		if err != nil {
			log.Panicf("cannot create the file: %s", err)
		}
		sw = provider.NewCsvWriter(csv.NewWriter(output))
	} else {
		sw = provider.NewWriter(os.Stdout)
		output = os.Stdout
	}

	defer output.Close()

	p, err := getProviderByName(*prPtr, *c)
	if err != nil {
		fmt.Printf("Can't reconize provider %s, %s\n", *prPtr, err)
		os.Exit(1)
	}

	err = run(p, sw, *n)
	if err != nil {
		log.Panic(err)
	}
}

func getProviderByName(name string, maxConcurecy uint) (provider.StoryProvider, error) {
	switch strings.ToLower(name) {
	case "hackernews":
		return hackernews.New(maxConcurecy), nil
	case "reddit":
		return reddit.New(), nil
	default:
		return nil, fmt.Errorf("provider not found")
	}
}

func run(p provider.StoryProvider, w provider.StoryWriter, n int) error {
	resp, err := p.GetStories(n)
	if err != nil {
		return fmt.Errorf("cannot get the stories: %s", err)
	}

	for r := range resp {
		if r.Error != nil {
			return r.Error
		}
		err := w.Write(r.Story)
		if err != nil {
			return fmt.Errorf("cannot write the story: %v", r.Error)
		}
	}

	return w.Flush()
}
