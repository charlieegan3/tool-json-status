package collect

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/pkg/errors"
)

type recentPlayResponse struct {
	RecentPlays []struct {
		Album     string `json:"Album"`
		Artist    string `json:"Artist"`
		Artwork   string `json:"Artwork"`
		Timestamp string `json:"Timestamp"`
		Track     string `json:"Track"`
	} `json:"RecentPlays"`
}

// LatestPlay holds data about the latest play from music.charlieegan3.com
type LatestPlay struct {
	Album   string `json:"album"`
	Artist  string `json:"artist"`
	Artwork string `json:"artwork"`
	Track   string `json:"track"`

	CreatedAt       time.Time `json:"created_at"`
	CreatedAtString string    `json:"created_at_string"`
}

// Collect returns a user's latest commit and project
func (l *LatestPlay) Collect(sourceURL string) error {
	resp, err := http.Get(fmt.Sprintf(sourceURL))
	if err != nil {
		return errors.Wrap(err, "Music get failed")
	}
	defer resp.Body.Close()

	var res recentPlayResponse
	err = json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		return errors.Wrap(err, "music body unmarshal failed")
	}

	if len(res.RecentPlays) < 1 {
		return errors.New("There were no plays in the recentPlayResponse")
	}

	latestPlay := res.RecentPlays[0]
	l.CreatedAt, err = time.Parse(time.RFC3339, latestPlay.Timestamp)
	if err != nil {
		return errors.Wrap(err, "Play time parsing failed")
	}

	l.Album = latestPlay.Album
	l.Artist = latestPlay.Artist
	l.Track = latestPlay.Track

	hostURL, _ := url.Parse(sourceURL)
	hostURL.Path = ""
	l.Artwork = hostURL.String() + latestPlay.Artwork

	return nil
}
