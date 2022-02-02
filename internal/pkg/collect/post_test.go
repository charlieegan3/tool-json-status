package collect

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPost(t *testing.T) {
	localServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, string(`{"location":"Hampstead Heath","url":"https://photos.charlieegan3.com/posts/2461","created_at":"2022-01-16T15:01:00Z"}`))
	}))

	var latestPost LatestPost
	err := latestPost.Collect(localServer.URL)
	if err != nil {
		t.Error(err)
	}

	if latestPost.Location != "Hampstead Heath" {
		t.Error(latestPost)
	}
}
