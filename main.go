package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/yanndr/topstories/provider"
	"github.com/yanndr/topstories/provider/reddit"
)

func main() {
	csvPtr := flag.Bool("csv", false, "Save the result to a csv file.")
	path := flag.String("o", "outupt.csv", "output file name")
	n := flag.Int("n", 20, "number of stories to display")
	// c := flag.Int("c", 20, "max concurency allowed")
	flag.Parse()

	cl := reddit.New()
	resp, err := cl.GetStories(*n)

	if err != nil {
		log.Panicf("cannot get the stories: %s", err)
	}

	if *csvPtr {
		f, err := os.Create(*path)
		if err != nil {
			log.Panicf("cannot create the file: %s", err)
		}
		defer f.Close()
		err = outputToCsv(f, resp)
		if err != nil {
			log.Panicf("cannot save to csv: %s", err)
		}
		return
	}

	err = handleResponse(resp, func(s provider.Story) error {
		_, err := fmt.Printf("|%-60s|%-100s|\n", s.Title(), s.URL())
		return err
	})
	if err != nil {
		log.Panic(err)
	}
}

func handleResponse(resp <-chan provider.Response, f func(provider.Story) error) error {
	for r := range resp {
		if r.Error != nil {
			return r.Error
		}
		err := f(r.Story)
		if err != nil {
			return err
		}
	}
	return nil
}

func outputToCsv(f io.Writer, resp <-chan provider.Response) error {
	w := csv.NewWriter(f)
	handleResponse(resp, func(s provider.Story) error {
		return w.Write([]string{s.Title(), s.URL()})
	})
	w.Flush()

	if err := w.Error(); err != nil {
		return fmt.Errorf("cannot write the csv file: %s", err)
	}

	return nil
}
