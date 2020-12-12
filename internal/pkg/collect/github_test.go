package collect

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGitHub(t *testing.T) {
	localServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		content, err := ioutil.ReadFile("github_response.json")
		if err != nil {
			t.Error(err)
		}
		fmt.Fprint(w, string(content))
	}))

	var latestCommit LatestCommit
	err := latestCommit.Collect(localServer.URL, "charlieegan3")

	if err != nil {
		t.Error(err)
	}
	if latestCommit.Repo.Name != "charlieegan3/dotfiles" {
		t.Error(latestCommit)
	}
	if latestCommit.Message != "first line" {
		t.Error(latestCommit)
	}
}
