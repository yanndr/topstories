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

type errorProvider struct {
}

func (errorProvider) GetStories(limit int) (<-chan provider.Response, error) {
	return nil, fmt.Errorf("error")
}

type fakeStoryErrorProvider struct {
}

func (fakeStoryErrorProvider) GetStories(limit int) (<-chan provider.Response, error) {
	ch := make(chan provider.Response)

	go func() {
		for i := 0; i < limit; i++ {
			ch <- provider.Response{Error: fmt.Errorf("error")}
		}
		close(ch)
	}()
	return ch, nil
}

type fakeStoryWriter struct {
	w io.Writer
}

func (w fakeStoryWriter) Write(s provider.Story) error {
	_, err := w.w.Write([]byte(fmt.Sprintf("|%-s|%s|\n", s.Title(), s.URL())))
	return err
}

func (fakeStoryWriter) Flush() error { return nil }

func TestRun(t *testing.T) {

	tt := []struct {
		name     string
		n        int
		provider provider.StoryProvider
		writter  provider.StoryWriter
		err      bool
	}{
		{"normal", 5, &fakeProvider{}, &fakeStoryWriter{}, false},
		{"provider error", 5, &errorProvider{}, &fakeStoryWriter{}, true},
		{"story error", 5, &fakeStoryErrorProvider{}, &fakeStoryWriter{}, true},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			b := &bytes.Buffer{}
			err := run(tc.provider, &fakeStoryWriter{w: b}, tc.n)

			if err != nil && !tc.err {
				t.Fatalf("expect no error, got %v", err)
			} else if err == nil && tc.err {
				t.Fatal("expect error, got no error")
			}

			if !tc.err && len(b.Bytes()) == 0 {
				t.Fatal("output should not be empty")
			}
		})
	}

}

func TestGetProviderByName(t *testing.T) {
	tt := []struct {
		name string
		err  bool
	}{
		{"HackerNews", false},
		{"hackernews", false},
		{"reddit", false},
		{"reDDit", false},
		{"false", true},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			p, err := getProviderByName(tc.name, 1)
			if err != nil && !tc.err {
				t.Fatalf("do not expect an error, got % v", err)
			} else if err == nil && tc.err {
				t.Fatal("expect an error, got no error")
			}
			if !tc.err {
				if p == nil {
					t.Fatalf("Expect a return value, got nil ")
				}
			}
		})
	}
}
