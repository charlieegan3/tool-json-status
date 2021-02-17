package collect

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/charlieegan3/json-charlieegan3/internal/pkg/proxy"
)

func TestInstagram(t *testing.T) {
	localServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var content []byte
		var err error
		if strings.Contains(r.URL.Path, "/p/") {
			content, err = ioutil.ReadFile("instagram_response_post.json")
		} else {
			content, err = ioutil.ReadFile("instagram_response_profile.json")
		}
		if err != nil {
			t.Error(err)
		}
		fmt.Fprint(w, string(content))
	}))

	localProxy := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		response, err := http.Get(q.Get("url"))

		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			fmt.Fprint(w, fmt.Sprintf("failed to ready body of downstream response: %s", err))
		}
		fmt.Fprint(w, string(body))
	}))

	proxy.Init(localProxy.URL, "MQ==")

	var latestPost LatestPost
	err := latestPost.Collect(localServer.URL, "charlieegan3", "MQ==")
	if err != nil {
		t.Error(err)
	}

	if latestPost.Location != "Dartmouth Park" {
		t.Error(latestPost)
	}

	if strings.Contains(latestPost.URL, "/p/CLXPx-rrl3P") == false {
		t.Error(latestPost)
	}

	if latestPost.Type != "photo" {
		t.Error(latestPost)
	}

	if fmt.Sprintf("%v", latestPost.CreatedAt) != "2021-02-16 18:31:14 +0000 GMT" {
		t.Error(latestPost.CreatedAt)
	}
}
