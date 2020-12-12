package collect

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestStrava(t *testing.T) {
	localServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		content, err := ioutil.ReadFile("strava_response.json")
		if err != nil {
			t.Error(err)
		}
		fmt.Fprint(w, string(content))
	}))

	var latestActivity LatestActivity
	err := latestActivity.Collect(localServer.URL)
	if err != nil {
		t.Error(err)
	}

	if latestActivity.Name != "Evaporate" {
		t.Error(latestActivity)
	}
	if latestActivity.Distance != 4231.5 {
		t.Error(latestActivity)
	}
	if latestActivity.MovingTime != 1470 {
		t.Error(latestActivity)
	}
	if latestActivity.AverageHeartrate != 142.9 {
		t.Error(latestActivity)
	}
	if latestActivity.Type != "Run" {
		t.Error(latestActivity)
	}
	if latestActivity.URL != "https://www.strava.com/activities/1748439744" {
		t.Error(latestActivity)
	}
}
