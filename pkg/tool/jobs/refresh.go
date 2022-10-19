package jobs

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/doug-martin/goqu/v9"

	"github.com/charlieegan3/tool-json-status/internal/pkg/status"
)

// Refresh will update the status in the database
type Refresh struct {
	DB *sql.DB

	ScheduleOverride string

	Username string

	PlaySource string
	PostSource string

	TwitterCredentials string
	StravaClientSecret string
	StravaRefreshToken string
	StravaClientID     string
}

func (r *Refresh) Name() string {
	return "refresh"
}

func (r *Refresh) Run(ctx context.Context) error {
	doneCh := make(chan bool)
	errCh := make(chan error)

	goquDB := goqu.New("postgres", r.DB)

	go func() {
		// get the current status
		sel := goquDB.From("jsonstatus.data").Select("value").Where(goqu.C("key").Eq("status")).Limit(1)
		var previousStatusJSON string
		found, err := sel.Executor().ScanVal(&previousStatusJSON)
		if err != nil {
			errCh <- fmt.Errorf("failed to get current status: %w", err)
			return
		}

		var previousStatus status.Data
		if found {
			err = json.Unmarshal([]byte(previousStatusJSON), &previousStatus)
			if err != nil {
				errCh <- fmt.Errorf("failed to unmarshal current status: %w", err)
				return
			}
		} else {
			previousStatus = status.Data{}
		}

		var s status.Data

		var wg sync.WaitGroup
		wg.Add(6)

		go func() {
			defer wg.Done()
			twitterCredentials := strings.Split(r.TwitterCredentials, ",")
			err := s.Tweet.Collect("https://api.twitter.com/1.1", twitterCredentials)
			if err != nil {
				log.Printf("failed to collect twitter: %s", err)
				s.Tweet = previousStatus.Tweet
			}
		}()

		go func() {
			defer wg.Done()
			err := s.Post.Collect(r.PostSource)
			if err != nil {
				log.Printf("failed to collect post: %s", err)
				s.Post = previousStatus.Post
			}
		}()

		go func() {
			defer wg.Done()
			err := s.Activity.Collect(r.StravaClientID, r.StravaClientSecret, r.StravaRefreshToken, "https://www.strava.com")
			if err != nil {
				log.Printf("failed to collect activity: %s", err)
				s.Activity = previousStatus.Activity
			}
		}()

		go func() {
			defer wg.Done()
			err := s.Film.Collect("https://letterboxd.com", r.Username)
			if err != nil {
				log.Printf("failed to collect film: %s", err)
				s.Film = previousStatus.Film
			}
		}()

		go func() {
			defer wg.Done()
			err := s.Commit.Collect("https://api.github.com", r.Username)
			if err != nil {
				log.Printf("failed to collect commit: %s", err)
				s.Commit = previousStatus.Commit
			}
		}()

		go func() {
			defer wg.Done()
			err := s.Play.Collect(r.PlaySource)
			if err != nil {
				log.Printf("failed to collect play: %s", err)
				s.Play = previousStatus.Play
			}
		}()

		wg.Wait()

		s.CreatedAt = time.Now()

		jsonStatus, err := json.Marshal(s)
		if err != nil {
			errCh <- fmt.Errorf("failed to marshal status: %w", err)
			return
		}

		// update the status
		ins := goquDB.Insert("jsonstatus.data").Rows(goqu.Record{"key": "status", "value": jsonStatus}).
			OnConflict(goqu.DoUpdate("key", goqu.Record{"value": jsonStatus}))
		_, err = ins.Executor().Exec()
		if err != nil {
			errCh <- fmt.Errorf("failed to update status: %w", err)
			return
		}

		doneCh <- true
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case e := <-errCh:
		return fmt.Errorf("job failed with error: %s", e)
	case <-doneCh:
		return nil
	}
}

func (r *Refresh) Timeout() time.Duration {
	return 15 * time.Second
}

func (r *Refresh) Schedule() string {
	if r.ScheduleOverride != "" {
		return r.ScheduleOverride
	}
	return "0 */5 * * * *"
}
