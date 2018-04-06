package json

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

// UnmarshalFromURL do a get request and parses the JSON-encoded data and stores the result
// in the value pointed to by v.
func UnmarshalFromURL(url string, v interface{}) error {
	c := &http.Client{}
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return fmt.Errorf("error creating the request: %s", err)
	}
	req.Header.Set("User-Agent", "github.com/yanndr/topStories")

	res, err := c.Do(req)
	if err != nil {
		return fmt.Errorf("error issuing the request: %s", err)
	}

	if 200 > res.StatusCode || res.StatusCode > 299 {
		return fmt.Errorf("err with response: %s", res.Status)
	}

	defer res.Body.Close()

	return parse(res.Body, v)
}

func parse(r io.Reader, v interface{}) error {
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return fmt.Errorf("parsing failed readall error: %v", err)
	}

	if len(b) == 0 {
		return nil
	}

	err = json.Unmarshal(b, v)
	if err != nil {
		return fmt.Errorf("parsing failed Unmasrshal error: %v", err)
	}

	return nil
}
