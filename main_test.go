package main

import (
	"bytes"
	"fmt"
	"io"
	"testing"

	"github.com/yanndr/topstories/provider"
)

type fakeStory struct {
}

func (fakeStory) Title() string {
	return "title"
}

func (fakeStory) URL() string {
	return "url"
}

type fakeProvider struct {
}

func (fakeProvider) GetStories(limit int) (<-chan provider.Response, error) {
	ch := make(chan provider.Response)

	go func() {
		for i := 0; i < limit; i++ {
			ch <- provider.Response{Story: &fakeStory{}}
		}
		close(ch)
	}()
	return ch, nil
}

type fakeStoryWriter struct {
	w io.Writer
}

func (w fakeStoryWriter) Write(s provider.Story) error {
	_, err := w.w.Write([]byte(fmt.Sprintf("|%-60s|%-100s|\n", s.Title(), s.URL())))
	return err
}

func (fakeStoryWriter) Flush() error { return nil }

func TestRun(t *testing.T) {

	tt := []struct {
		name    string
		story   fakeStory
		errItem error
		errFunc error
	}{
		{name: "Response with story", story: fakeStory{}},
		{name: "Response with error", errItem: fmt.Errorf("error"), errFunc: nil},
		{name: "Response with func error", errFunc: fmt.Errorf("error")},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			b := &bytes.Buffer{}
			run(&fakeProvider{}, &fakeStoryWriter{w: b}, 10)
		})
	}
}
