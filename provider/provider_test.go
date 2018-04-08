package provider

import (
	"bytes"
	"encoding/csv"
	"testing"
)

type fakeStory struct {
}

func (fakeStory) Title() string {
	return "title"
}

func (fakeStory) URL() string {
	return "url"
}

func TestNewWriter(t *testing.T) {
	b := &bytes.Buffer{}
	w := NewWriter(b)

	w.Write(&fakeStory{})

	if len(b.Bytes()) == 0 {
		t.Fatal("Expected data written on the buffer, got nothig")
	}
}

func TestNewCsvWriter(t *testing.T) {
	b := &bytes.Buffer{}
	w := NewCsvWriter(csv.NewWriter(b))

	w.Write(&fakeStory{})
	w.Flush()

	if len(b.Bytes()) == 0 {
		t.Fatal("Expected data written on the buffer, got nothig")
	}
}
