package tool

import (
	"database/sql"
	"embed"
	"fmt"

	"github.com/Jeffail/gabs/v2"
	"github.com/charlieegan3/toolbelt/pkg/apis"
	"github.com/gorilla/mux"

	"github.com/charlieegan3/tool-json-status/pkg/tool/handlers"
	"github.com/charlieegan3/tool-json-status/pkg/tool/jobs"
)

//go:embed migrations
var jsonStatusToolMigrations embed.FS

// JSONStatus is a tool for generating a personal JSON status for my public activities
type JSONStatus struct {
	config *gabs.Container
	db     *sql.DB
}

func (t *JSONStatus) Name() string {
	return "json-status"
}

func (t *JSONStatus) FeatureSet() apis.FeatureSet {
	return apis.FeatureSet{
		Config:   true,
		Jobs:     true,
		Database: true,
		HTTP:     true,
	}
}

func (t *JSONStatus) SetConfig(config map[string]any) error {
	t.config = gabs.Wrap(config)

	return nil
}
func (t *JSONStatus) Jobs() ([]apis.Job, error) {
	var j []apis.Job
	var path string
	var ok bool

	// load all config
	path = "jobs.refresh.schedule"
	scheduleRefresh, ok := t.config.Path(path).Data().(string)
	if !ok {
		return j, fmt.Errorf("missing required config path: %s", path)
	}

	path = "jobs.check.schedule"
	scheduleCheck, ok := t.config.Path(path).Data().(string)
	if !ok {
		return j, fmt.Errorf("missing required config path: %s", path)
	}

	path = "jobs.check.alert_endpoint"
	alertEndpoint, ok := t.config.Path(path).Data().(string)
	if !ok {
		return j, fmt.Errorf("missing required config path: %s", path)
	}

	path = "username"
	username, ok := t.config.Path(path).Data().(string)
	if !ok {
		return j, fmt.Errorf("missing required config path: %s", path)
	}

	path = "twitter.credentials"
	twitterCredentials, ok := t.config.Path(path).Data().(string)
	if !ok {
		return j, fmt.Errorf("missing required config path: %s", path)
	}

	path = "strava.client_secret"
	stravaClientSecret, ok := t.config.Path(path).Data().(string)
	if !ok {
		return j, fmt.Errorf("missing required config path: %s", path)
	}

	path = "strava.refresh_token"
	stravaRefreshToken, ok := t.config.Path(path).Data().(string)
	if !ok {
		return j, fmt.Errorf("missing required config path: %s", path)
	}

	path = "strava.client_id"
	stravaClientID, ok := t.config.Path(path).Data().(string)
	if !ok {
		return j, fmt.Errorf("missing required config path: %s", path)
	}

	path = "play_source"
	playSource, ok := t.config.Path(path).Data().(string)
	if !ok {
		return j, fmt.Errorf("missing required config path: %s", path)
	}
	path = "post_source"
	postSource, ok := t.config.Path(path).Data().(string)
	if !ok {
		return j, fmt.Errorf("missing required config path: %s", path)
	}

	return []apis.Job{
		&jobs.Check{
			DB:               t.db,
			ScheduleOverride: scheduleCheck,
			AlertEndpoint:    alertEndpoint,
		},
		&jobs.Refresh{
			DB:                 t.db,
			ScheduleOverride:   scheduleRefresh,
			Username:           username,
			PlaySource:         playSource,
			PostSource:         postSource,
			TwitterCredentials: twitterCredentials,
			StravaClientSecret: stravaClientSecret,
			StravaRefreshToken: stravaRefreshToken,
			StravaClientID:     stravaClientID,
		},
	}, nil
}
func (t *JSONStatus) ExternalJobsFuncSet(f func(job apis.ExternalJob) error) {}

func (t *JSONStatus) DatabaseMigrations() (*embed.FS, string, error) {
	return &jsonStatusToolMigrations, "migrations", nil
}
func (t *JSONStatus) DatabaseSet(db *sql.DB) {
	t.db = db
}

func (t *JSONStatus) HTTPPath() string { return "json-status" }
func (t *JSONStatus) HTTPAttach(router *mux.Router) error {
	router.HandleFunc(
		"/latest.json",
		handlers.BuildLatestHandler(t.db),
	).Methods("GET", "OPTIONS")

	return nil
}
