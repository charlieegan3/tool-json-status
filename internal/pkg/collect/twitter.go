package collect

import (
	"fmt"
	"net/url"
	"time"

	"github.com/ChimeraCoder/anaconda"
)

// LatestTweet wraps the required data for a tweet
type LatestTweet struct {
	Text            string    `json:"text"`
	Link            string    `json:"link"`
	Location        string    `json:"location"`
	CreatedAt       time.Time `json:"created_at"`
	CreatedAtString string    `json:"created_at_string"`
}

// Collect returns the latest tweet for the requesting user
// Use https://api.twitter.com/1.1 as the baseURL
func (l *LatestTweet) Collect(baseURL string, credentials []string) error {
	api := anaconda.NewTwitterApiWithCredentials(credentials[0], credentials[1], credentials[2], credentials[3])
	api.SetBaseUrl(baseURL)

	params := url.Values{}
	params.Set("include_entities", "false")
	data, err := api.GetUserTimeline(params)
	if err != nil {
		return err
	}

	createdAt, err := data[0].CreatedAtTime()
	if err != nil {
		return err
	}

	l.Text = data[0].Text
	l.CreatedAt = createdAt
	l.Location = data[0].Place.Name
	l.Link = fmt.Sprintf("https://twitter.com/%s/status/%s", data[0].User.ScreenName, data[0].IdStr)

	return nil
}
