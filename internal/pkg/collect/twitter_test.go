package collect

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestTwitter(t *testing.T) {
	localServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		content, err := ioutil.ReadFile("twitter_response.json")
		if err != nil {
			t.Error(err)
		}
		fmt.Fprint(w, string(content))
	}))

	var latestTweet LatestTweet
	err := latestTweet.Collect(localServer.URL, []string{"t", "t", "t", "t"})
	if err != nil {
		t.Error(err)
	}

	if latestTweet.Text != "just another test" {
		t.Error(latestTweet.Text)
	}
	if latestTweet.Link != "https://twitter.com/oauth_dancer/status/240558470661799936" {
		t.Error(latestTweet.Link)
	}
	if fmt.Sprintf("%v", latestTweet.CreatedAt) != "2012-08-28 21:16:23 +0000 +0000" {
		t.Error(latestTweet.CreatedAt)
	}
	if latestTweet.Location != "Berkeley" {
		t.Error(latestTweet.Location)
	}
}
