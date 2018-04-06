package main

import (
	"bytes"
	"fmt"
	"log"
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
func TestHandleResponse(t *testing.T) {

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
			c := make(chan provider.Response)
			go func() {
				defer close(c)
				c <- provider.Response{Story: tc.story, Error: tc.errItem}
			}()

			called := false
			err := handleResponse(c, func(s provider.Story) error {
				called = true
				return tc.errFunc
			})

			if err != tc.errFunc && err != tc.errItem {
				t.Fatalf("Expected %v got %v", tc.errFunc, err)
			}

			if tc.errItem != nil && called {
				t.Fatal("Function not expected to be call")
			} else if tc.errItem == nil && !called {
				t.Fatal("Function  expected to be call")
			}
		})
	}
}

func TestOutputToCsv(t *testing.T) {
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
			c := make(chan provider.Response)
			go func() {
				defer close(c)
				c <- provider.Response{Story: tc.story, Error: tc.errItem}
			}()

			buf := bytes.NewBufferString("")
			outputToCsv(buf, c)
			if len(buf.String()) == 0 && tc.errItem == nil && tc.errFunc == nil {
				log.Fatalf("Expected value got no value on the buffer.")
			}
		})
	}
}
