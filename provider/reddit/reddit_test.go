package reddit

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
func TestGetStories(t *testing.T) {

	t.Parallel()

	tt := []struct {
		name     string
		n        int
		response string
		err      bool
	}{
		{name: "happy", n: 1, response: response, err: false},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			mux := http.NewServeMux()
			srv := httptest.NewServer(mux)

			c := &reddit{
				newsURL: srv.URL,
			}

			mux.HandleFunc("/", handleString(tc.response))

			resp, err := c.GetStories(tc.n)

			if (err != nil) != tc.err {
				if tc.err {
					t.Fatalf("expected error, got no error")
				}
				t.Fatalf("expected no error, got %v", err)
			}
			if err != nil && tc.err {
				return
			}

			i := 0
			var errs []error
			for r := range resp {
				i++
				if r.Error != nil {
					errs = append(errs, r.Error)
					continue
				}
				if r.Story == nil {
					t.Fatal("error Story not espected to be nil")
				}
				if r.Story.Title() != "Test" {
					t.Fatalf("error expected test %v", r.Story.Title())
				}

				if r.Story.URL() != "https://test.html" {
					t.Fatalf("error expected https://test.html got %v", r.Story.URL())
				}
			}

			if i != tc.n {
				t.Fatalf("error, expect %v got %v", tc.n, i)
			}
		})
	}
}

const response string = `{"kind": "Listing", "data": {"after": "t3_8ad4v0", "dist": 1, "modhash": "", "whitelist_status": "all_ads", "children": [{"kind": "t3", "data": {"subreddit_id": "t5_2rc7j", "approved_at_utc": null, "send_replies": true, "mod_reason_by": null, "banned_by": null, "num_reports": null, "removal_reason": null, "subreddit": "golang", "selftext_html": null, "selftext": "", "likes": null, "suggested_sort": null, "user_reports": [], "secure_media": null, "is_reddit_media_domain": true, "saved": false, "id": "8ad4v0", "banned_at_utc": null, "mod_reason_title": null, "view_count": null, "archived": false, "clicked": false, "no_follow": true, "author": "kaveman98", "num_crossposts": 0, "link_flair_text": null, "mod_reports": [], "can_mod_post": false, "is_crosspostable": false, "pinned": false, "score": 2, "approved_by": null, "over_18": false, "report_reasons": null, "domain": "i.redd.it", "hidden": false, "thumbnail": "", "edited": false, "link_flair_css_class": null, "author_flair_css_class": null, "contest_mode": false, "gilded": 0, "downs": 0, "brand_safe": true, "secure_media_embed": {}, "media_embed": {}, "author_flair_text": null, "stickied": false, "visited": false, "can_gild": false, "is_self": false, "parent_whitelist_status": "all_ads", "name": "t3_8ad4v0", "spoiler": false, "permalink": "/r/golang/comments/8ad4v0/is_there_a_better_way_to_do_the_following_ty_i_am/", "subreddit_type": "public", "locked": false, "hide_score": false, "created": 1523077824.0, "url": "https://test.html", "whitelist_status": "all_ads", "quarantine": false, "subreddit_subscribers": 44134, "created_utc": 1523049024.0, "subreddit_name_prefixed": "r/golang", "ups": 2, "media": null, "num_comments": 4, "title": "Test", "mod_note": null, "is_video": false, "distinguished": null}}], "before": null}}`
