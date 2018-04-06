package json

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func handleString(s string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, s)
	}
}

func TestUnmarshalFromURL(t *testing.T) {
	tt := []struct {
		name   string
		input  string
		result []int
		err    bool
	}{
		{name: "basic", input: "[1,2,3,4]", result: []int{1, 2, 3, 4}, err: false},
		{name: "empty", input: "", result: []int{}, err: false},
		{name: "bignumber", input: "[16755530,16760736,16756901,16761349,16757044]", result: []int{16755530, 16760736, 16756901, 16761349, 16757044}, err: false},
		{name: "wrong format", input: "test", result: nil, err: true},
		{name: "wrong number", input: "[0,1,01,1]", result: nil, err: true},
		{name: "with letter", input: "[0,1,b,1]", result: nil, err: true},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			mux := http.NewServeMux()
			srv := httptest.NewServer(mux)
			mux.HandleFunc("/", handleString(tc.input))

			var output []int
			UnmarshalFromURL(srv.URL, &output)
			if len(output) != len(tc.result) {
				t.Fatalf("expected %v, got %v", len(tc.result), len(output))
			}

			for i := 0; i < len(output); i++ {
				if output[i] != tc.result[i] {
					t.Fatalf("expected %v, got %v", tc.result[i], output[i])
				}
			}
		})
	}
}
