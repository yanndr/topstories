package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/yanndr/topstories/provider"
	"github.com/yanndr/topstories/provider/hackernews"
)

func main() {
	var (
		output *os.File
		sw     provider.StoryWriter
		err    error
	)

	csvPtr := flag.Bool("csv", false, "Save the result to a csv file.")
	path := flag.String("o", "outupt.csv", "output file name")
	n := flag.Int("n", 20, "number of stories to display")
	c := flag.Int("c", 20, "max concurency allowed")
	flag.Parse()

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
	p := hackernews.New(*c)
	err = run(p, sw, *n)
	if err != nil {
		log.Panic(err)
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
			return r.Error
		}
	}

	if err != nil {
		return err
	}
	return w.Flush()
}
