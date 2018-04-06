package provider

import (
	"encoding/csv"
	"fmt"
	"io"
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

type StoryWriter interface {
	Write(s Story) error
	Flush() error
}

type Writer struct {
	w io.Writer
}

func NewWriter(w io.Writer) *Writer {
	return &Writer{
		w: w,
	}
}

func (w Writer) Write(s Story) error {
	_, err := w.w.Write([]byte(fmt.Sprintf("|%-60s|%-100s|\n", s.Title(), s.URL())))
	return err
}

func (Writer) Flush() error { return nil }

type CsvWriter struct {
	w *csv.Writer
}

func NewCsvWriter(w *csv.Writer) *CsvWriter {
	return &CsvWriter{
		w: w,
	}
}

func (w CsvWriter) Write(s Story) error {
	return w.w.Write([]string{s.Title(), s.URL()})
}

func (w CsvWriter) Flush() error {
	w.w.Flush()
	return w.w.Error()
}
