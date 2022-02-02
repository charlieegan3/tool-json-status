package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	backoff "github.com/cenkalti/backoff/v4"
	"github.com/dustin/go-humanize"
	"github.com/go-yaml/yaml"
	"github.com/pkg/errors"

	"github.com/charlieegan3/json-charlieegan3/internal/pkg/collect"
	"github.com/charlieegan3/json-charlieegan3/internal/pkg/proxy"
)

type config struct {
	ExportPath         string `yaml:"export_path"`
	LastfmKey          string `yaml:"lastfm_key"`
	PlaySource         string `yaml:"play_source"`
	ProxyToken         string `yaml:"proxy_token"`
	ProxyURL           string `yaml:"proxy_url"`
	StatusHost         string `yaml:"status_host"`
	StatusKey          string `yaml:"status_key"`
	StravaClientID     string `yaml:"strava_client_id"`
	StravaClientSecret string `yaml:"strava_client_secret"`
	StravaRefreshToken string `yaml:"strava_refresh_token"`
	TwitterCredentials string `yaml:"twitter_credentials"`
	Username           string `yaml:"username"`
}

var cfg config

type status struct {
	Tweet    collect.LatestTweet    `json:"tweet"`
	Post     collect.LatestPost     `json:"post"`
	Activity collect.LatestActivity `json:"activity"`
	Film     collect.LatestFilm     `json:"film"`
	Commit   collect.LatestCommit   `json:"commit"`
	Play     collect.LatestPlay     `json:"play"`
}

func (s *status) setCreatedAtStrings() {
	s.Tweet.CreatedAtString = compactHumanizeTime(s.Tweet.CreatedAt)
	s.Post.CreatedAtString = compactHumanizeTime(s.Post.CreatedAt)
	s.Activity.CreatedAtString = compactHumanizeTime(s.Activity.CreatedAt)
	s.Film.CreatedAtString = compactHumanizeTime(s.Film.CreatedAt)
	s.Commit.CreatedAtString = compactHumanizeTime(s.Commit.CreatedAt)
	s.Play.CreatedAtString = compactHumanizeTime(s.Play.CreatedAt)
}

func (s *status) fetchCurrent(url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return errors.Wrap(err, "current status get failed")
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&s)
	if err != nil {
		return errors.Wrap(err, "current status unmarshal failed")
	}

	return nil
}

func (s *status) fetchNew(previousStatus status) {
	username := cfg.Username

	var wg sync.WaitGroup
	wg.Add(6)

	go func() {
		defer wg.Done()
		twitterCredentials := strings.Split(cfg.TwitterCredentials, ",")
		err := s.Tweet.Collect("https://api.twitter.com/1.1", twitterCredentials)
		if err != nil {
			fmt.Println(errors.Wrap(err, "twitter error"))
			s.Tweet = previousStatus.Tweet
		}
	}()

	go func() {
		defer wg.Done()
		err := s.Post.Collect("https://photos.charlieegan3.com/posts/latest.json")
		if err != nil {
			fmt.Println(errors.Wrap(err, "post collection error"))
			s.Post = previousStatus.Post
		}
	}()

	go func() {
		defer wg.Done()
		err := s.Activity.Collect("https://www.strava.com")
		if err != nil {
			fmt.Println(errors.Wrap(err, "strava error"))
			s.Activity = previousStatus.Activity
		}
	}()

	go func() {
		defer wg.Done()
		err := s.Film.Collect("https://letterboxd.com", username)
		if err != nil {
			fmt.Println(errors.Wrap(err, "letterboxd error"))
			s.Film = previousStatus.Film
		}
	}()

	go func() {
		defer wg.Done()
		err := s.Commit.Collect("https://api.github.com", username)
		if err != nil {
			fmt.Println(errors.Wrap(err, "github error"))
			s.Commit = previousStatus.Commit
		}
	}()

	go func() {
		defer wg.Done()
		err := s.Play.Collect(cfg.PlaySource)
		if err != nil {
			fmt.Println(errors.Wrap(err, "plays error"))
			s.Play = previousStatus.Play
		}
	}()

	wg.Wait()
}

func compactHumanizeTime(time time.Time) string {
	humanReadable := humanize.Time(time)

	humanReadable = strings.Replace(humanReadable, " year", "yr", 1)
	humanReadable = strings.Replace(humanReadable, " month", "mth", 1)
	humanReadable = strings.Replace(humanReadable, " weeks", "w", 1)
	humanReadable = strings.Replace(humanReadable, " week", "w", 1)
	humanReadable = strings.Replace(humanReadable, " days", "d", 1)
	humanReadable = strings.Replace(humanReadable, " day", "d", 1)
	humanReadable = strings.Replace(humanReadable, " hours", "h", 1)
	humanReadable = strings.Replace(humanReadable, " hour", "h", 1)
	humanReadable = strings.Replace(humanReadable, " minutes", "m", 1)
	humanReadable = strings.Replace(humanReadable, " minute", "m", 1)
	humanReadable = strings.Replace(humanReadable, " seconds", "s", 1)
	humanReadable = strings.Replace(humanReadable, " second", "s", 1)

	return humanReadable
}

func writeFile(content []byte, path string) error {
	file, err := os.Create(path)
	if err != nil {
		return errors.Wrap(err, "failed to open status file")
	}
	defer file.Close()

	_, err = fmt.Fprint(file, string(content))
	if err != nil {
		return errors.Wrap(err, "failed to write status file")
	}

	return nil
}

func downloadStatusData(existingURL string, refresh bool, output string) error {
	var previousStatus status
	var nextStatus status

	log.Println("fetching previous data")
	err := previousStatus.fetchCurrent(existingURL)
	if err != nil {
		return errors.Wrap(err, "error getting current status")
	}

	if refresh {
		log.Println("fetching latest data")
		nextStatus.fetchNew(previousStatus)
	} else {
		log.Println("re-using previous data")
		nextStatus = previousStatus
	}

	nextStatus.setCreatedAtStrings()

	if nextStatus == previousStatus {
		log.Println("no update required, exiting")
		return nil
	}

	log.Println("formatting data for export")
	statusJSON, err := json.Marshal(nextStatus)
	if err != nil {
		return errors.Wrap(err, "json generation error")
	}

	log.Println("saving data")
	err = writeFile(statusJSON, output)
	if err != nil {
		return errors.Wrap(err, "failed to write file")
	}

	log.Println("completed successfully")
	return nil
}

func loadConfig(configPath string) error {
	file, err := os.Open(configPath)
	if err != nil {
		return errors.Wrap(err, "failed to open file")
	}
	defer file.Close()

	yamlBlob, err := ioutil.ReadAll(file)
	if err != nil {
		return errors.Wrap(err, "failed to read file")
	}

	err = yaml.Unmarshal(yamlBlob, &cfg)
	if err != nil {
		return errors.Wrap(err, "failed to parse yaml")
	}

	return nil
}

func main() {
	output := flag.String("output", "status.json", "where to save the json output")
	existingURL := flag.String("existing-url", "https://charlieegan3.github.io/json-charlieegan3/build/status.json", "where to get the existing status from")
	refresh := flag.Bool("refresh", false, "download new status data")
	repeat := flag.Bool("repeat", false, "run forever refreshing the data every interval")
	interval := flag.Int("interval", 600, "seconds to wait before refreshes")
	configPath := flag.String("config", "config.yaml", "where to load config from")
	flag.Parse()

	// load the config file
	if err := loadConfig(*configPath); err != nil {
		log.Fatal(err)
	}

	// init the proxy
	if err := proxy.Init(cfg.ProxyURL, cfg.ProxyToken); err != nil {
		log.Fatal(err)
	}
	// init strava
	if err := collect.StravaInit(cfg.StravaClientID, cfg.StravaClientSecret, cfg.StravaRefreshToken); err != nil {
		log.Fatal(err)
	}

	if *repeat {
		log.Println("starting loop")

		for {
			downloadStatusData(*existingURL, *refresh, *output)
			log.Printf("sleeping for %d", *interval)
			time.Sleep(time.Duration(*interval) * time.Second)
		}
	} else {
		// if running as a one off, then retry the job for 2 mins
		b := backoff.NewExponentialBackOff()
		b.MaxElapsedTime = 2 * time.Minute

		err := backoff.Retry(func() error {
			return downloadStatusData(*existingURL, *refresh, *output)
		}, b)

		if err != nil {
			log.Fatalf("failed to refresh data after retry: %s", err)
		}
	}
}
