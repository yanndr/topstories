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

// StoryWriter is an interface that define the methods of a story writer.
type StoryWriter interface {
	Write(s Story) error
	Flush() error
}

type writer struct {
	w io.Writer
}

// NewWriter returns a new story writer to write on w.
func NewWriter(w io.Writer) StoryWriter {
	return &writer{
		w: w,
	}
}

func (w *writer) Write(s Story) error {
	_, err := w.w.Write([]byte(fmt.Sprintf("|%-70s|%-100s|\n", s.Title(), s.URL())))
	return err
}

func (*writer) Flush() error { return nil }

type csvWriter struct {
	w *csv.Writer
}

// NewCsvWriter returns a new story CSV writer to write on w.
func NewCsvWriter(w *csv.Writer) StoryWriter {
	return &csvWriter{
		w: w,
	}
}

func (w *csvWriter) Write(s Story) error {
	return w.w.Write([]string{s.Title(), s.URL()})
}

func (w *csvWriter) Flush() error {
	w.w.Flush()
	return w.w.Error()
}
