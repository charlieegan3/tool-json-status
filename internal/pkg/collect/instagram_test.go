package collect

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestInstagram(t *testing.T) {
	localServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var content []byte
		var err error
		if strings.Contains(r.URL.Path, "/p/") {
			content, err = ioutil.ReadFile("instagram_response_post.html")
		} else {
			content, err = ioutil.ReadFile("instagram_response_profile.html")
		}
		if err != nil {
			t.Error(err)
		}
		fmt.Fprint(w, string(content))
	}))

	var latestPost LatestPost
	err := latestPost.Collect(localServer.URL, "charlieegan3", "MQo=")
	if err != nil {
		t.Error(err)
	}

	if latestPost.Location != "Barbican Estate" {
		t.Error(latestPost)
	}

	if strings.Contains(latestPost.URL, "/p/BmCO0mAgC2h") == false {
		t.Error(latestPost)
	}

	if latestPost.Type != "photo" {
		t.Error(latestPost)
	}

	if fmt.Sprintf("%v", latestPost.CreatedAt) != "2018-08-04 00:17:13 +0100 BST" {
		t.Error(latestPost.CreatedAt)
	}
}
