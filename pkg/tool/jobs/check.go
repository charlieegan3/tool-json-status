package jobs

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/doug-martin/goqu/v9"

	"github.com/charlieegan3/tool-json-status/internal/pkg/status"
)

// Check will make sure that the status data is fresh and alert if not
type Check struct {
	DB *sql.DB

	ScheduleOverride string

	AlertEndpoint string
}

func (c *Check) Name() string {
	return "check"
}

func (c *Check) Run(ctx context.Context) error {
	doneCh := make(chan bool)
	errCh := make(chan error)

	goquDB := goqu.New("postgres", c.DB)

	go func() {
		var err error

		defer func() {
			if err != nil {
				requestError := alert(c.AlertEndpoint, "json-status: Check Error", err.Error())
				if requestError != nil {
					errCh <- fmt.Errorf("failed to alert on error %s: %w", err.Error(), requestError)
				}
				errCh <- err
			}
		}()

		// get the current status
		sel := goquDB.From("jsonstatus.data").Select("value").Where(goqu.C("key").Eq("status")).Limit(1)
		var latestStatusJSON string
		var found bool
		found, err = sel.Executor().ScanVal(&latestStatusJSON)
		if err != nil {
			err = fmt.Errorf("failed to get current status: %w", err)
			return
		}

		var latestStatus status.Data
		if found {
			err = json.Unmarshal([]byte(latestStatusJSON), &latestStatus)
			if err != nil {
				err = fmt.Errorf("failed to unmarshal current status: %w", err)
				return
			}
		} else {
			latestStatus = status.Data{}
		}

		if latestStatus.CreatedAt.Sub(time.Now()) > time.Hour {
			err = fmt.Errorf("status is stale, last updated: %s", latestStatus.CreatedAt)
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

func (c *Check) Timeout() time.Duration {
	return 15 * time.Second
}

func (c *Check) Schedule() string {
	if c.ScheduleOverride != "" {
		return c.ScheduleOverride
	}
	return "0 */30 * * * *"
}

func alert(webhookRSSEndpoint, title, message string) error {
	datab := []map[string]string{
		{
			"title": title,
			"body":  message,
			"url":   "",
		},
	}

	b, err := json.Marshal(datab)
	if err != nil {
		return fmt.Errorf("failed to form alert item JSON: %s", err)
	}

	client := &http.Client{}
	req, err := http.NewRequest("POST", webhookRSSEndpoint, bytes.NewBuffer(b))
	if err != nil {
		return fmt.Errorf("failed to build request for alert item: %s", err)
	}

	req.Header.Add("Content-Type", "application/json; charset=utf-8")

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request for alert item: %s", err)
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to send request: non 200OK response")
	}

	return nil
}
