package collect

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestLetterboxd(t *testing.T) {
	localServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		content, err := ioutil.ReadFile("letterboxd_response.rss")
		if err != nil {
			t.Error(err)
		}
		fmt.Fprint(w, string(content))
	}))

	var latestFilm LatestFilm
	err := latestFilm.Collect(localServer.URL, "charlieegan3")
	if err != nil {
		t.Error(err)
	}

	if latestFilm.Title != "The Train Stop" {
		t.Error(latestFilm)
	}
	if latestFilm.Year != "2000" {
		t.Error(latestFilm)
	}
	if !strings.Contains(fmt.Sprintf("%v", latestFilm.CreatedAt), "2020-06-14 08:01:06") {
		t.Errorf("%v", latestFilm.CreatedAt)
	}
	if latestFilm.Rating != "" {
		t.Error(latestFilm.Rating)
	}
	if latestFilm.Link != "https://letterboxd.com/charlieegan3/film/the-train-stop/" {
		t.Error(latestFilm.Link)
	}
}
