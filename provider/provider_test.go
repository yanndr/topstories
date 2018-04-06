package provider

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetResponseBody(t *testing.T) {
	const response string = "response"
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, response)
	}))

	body, _ := GetResponseBody(srv.URL)

	defer body.Close()

	b, _ := ioutil.ReadAll(body)
	s := string(b)

	if s != response {
		t.Fatalf("Error expected %s, got %s", response, s)
	}
}
