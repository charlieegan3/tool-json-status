package collect

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMusic(t *testing.T) {
	localServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		content, err := ioutil.ReadFile("music_response.json")
		if err != nil {
			t.Error(err)
		}
		fmt.Fprint(w, string(content))
	}))

	var latestPlay LatestPlay
	err := latestPlay.Collect(localServer.URL)

	if err != nil {
		t.Error(err)
	}
	if latestPlay.Track != "Phase" {
		t.Error(latestPlay)
	}
	if latestPlay.Artist != "Krrum" {
		t.Error(latestPlay)
	}
}
