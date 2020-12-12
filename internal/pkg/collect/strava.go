package collect

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/pkg/errors"
)

var stravaClientID, stravaClientSecret, stravaRefreshToken string

// StravaInit sets the credentials for the Collect function to talk to the
// strava api
func StravaInit(clientID, clientSecret, refreshToken string) error {
	if clientID == "" {
		return fmt.Errorf("client id cannot be blank")
	}
	if clientSecret == "" {
		return fmt.Errorf("client secret cannot be blank")
	}
	if refreshToken == "" {
		return fmt.Errorf("refresh token cannot be blank")
	}

	stravaClientID = clientID
	stravaClientSecret = clientSecret
	stravaRefreshToken = refreshToken

	return nil
}

type activity struct {
	AverageHeartrate float64 `json:"average_heartrate"`
	ID               int64   `json:"id"`
	Distance         float64 `json:"distance"`
	MovingTime       int64   `json:"moving_time"`
	Name             string  `json:"name"`
	StartDate        string  `json:"start_date"`
	Type             string  `json:"type"`
}

// LatestActivity wraps deta about the latest activity
type LatestActivity struct {
	AverageHeartrate float64   `json:"average_heartrate"`
	URL              string    `json:"url"`
	Distance         float64   `json:"distance"`
	MovingTime       int64     `json:"moving_time"`
	Name             string    `json:"name"`
	Type             string    `json:"type"`
	CreatedAt        time.Time `json:"created_at"`
	CreatedAtString  string    `json:"created_at_string"`
}

type accessTokenResponse struct {
	AccessToken  string `json:"access_token"`
	ExpiresAt    int64  `json:"expires_at"`
	ExpiresIn    int64  `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
}

// Collect returns details about the most recent strava activity
// host https://www.strava.com
func (l *LatestActivity) Collect(host string) error {
	body := strings.NewReader(
		fmt.Sprintf(
			`client_id=%s&client_secret=%s&grant_type=refresh_token&refresh_token=%s`,
			stravaClientID,
			stravaClientSecret,
			stravaRefreshToken))

	req, err := http.NewRequest("POST", "https://www.strava.com/api/v3/oauth/token", body)
	if err != nil {
		return errors.Wrap(err, "failed to build strava access token request")
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return errors.Wrap(err, "failed to get strava access token")
	}
	defer resp.Body.Close()

	var tokenResponse accessTokenResponse
	err = json.NewDecoder(resp.Body).Decode(&tokenResponse)
	if err != nil {
		return errors.Wrap(err, "access token body unmarshal failed")
	}

	req, err = http.NewRequest("GET", fmt.Sprintf("%s/api/v3/athlete/activities", host), nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", tokenResponse.AccessToken))
	if err != nil {
		return errors.Wrap(err, "build request failed")
	}
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		return errors.Wrap(err, "get activities failed")
	}

	defer resp.Body.Close()

	var activities []activity
	err = json.NewDecoder(resp.Body).Decode(&activities)
	if err != nil {
		return errors.Wrap(err, "body unmarshal failed")
	}

	if len(activities) == 0 {
		return errors.New("no activities found")
	}

	activity := activities[0]
	createdAt, err := time.Parse(time.RFC3339, activity.StartDate)
	if err != nil {
		return errors.Wrap(err, "latest activity time parsing failed")
	}

	l.AverageHeartrate = activity.AverageHeartrate
	l.Distance = activity.Distance
	l.MovingTime = activity.MovingTime
	l.Name = activity.Name
	l.Type = activity.Type
	l.CreatedAt = createdAt
	l.URL = fmt.Sprintf("https://www.strava.com/activities/%d", activity.ID)

	return nil
}
