package collect

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

// LatestPost stores the URL, location and time of the latest post, video or photo
type LatestPost struct {
	URL             string    `json:"url"`
	Location        string    `json:"location"`
	CreatedAt       time.Time `json:"created_at"`
	CreatedAtString string    `json:"created_at_string"`
}

// Collect returns latest post
func (l *LatestPost) Collect(url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to get latest post data: %s", err)
	}

	body, readErr := ioutil.ReadAll(resp.Body)
	if readErr != nil {
		return fmt.Errorf("failed to read repsonse: %s", err)
	}

	jsonErr := json.Unmarshal(body, l)
	if jsonErr != nil {
		return fmt.Errorf("failed to unmarshal as post: %s", err)
	}

	return nil
}
