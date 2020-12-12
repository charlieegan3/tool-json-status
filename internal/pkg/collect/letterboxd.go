package collect

import (
	"encoding/xml"
	"fmt"
	"net/http"
	"regexp"
	"time"

	"github.com/pkg/errors"
)

type rssDocument struct {
	XMLName xml.Name `xml:"rss,omitempty" json:"rss,omitempty"`
	Channel *struct {
		XMLName xml.Name `xml:"channel,omitempty" json:"channel,omitempty"`
		Item    []*struct {
			XMLName xml.Name `xml:"item,omitempty" json:"item,omitempty"`
			Link    *struct {
				XMLName xml.Name `xml:"link,omitempty" json:"link,omitempty"`
				String  string   `xml:",chardata" json:",omitempty"`
			} `xml:"link,omitempty" json:"link,omitempty"`
			PubDate *struct {
				XMLName xml.Name `xml:"pubDate,omitempty" json:"pubDate,omitempty"`
				String  string   `xml:",chardata" json:",omitempty"`
			} `xml:"pubDate,omitempty" json:"pubDate,omitempty"`
			Title *struct {
				XMLName xml.Name `xml:"title,omitempty" json:"title,omitempty"`
				String  string   `xml:",chardata" json:",omitempty"`
			} `xml:"title,omitempty" json:"title,omitempty"`
		} `xml:"item,omitempty" json:"item,omitempty"`
	} `xml:"channel,omitempty" json:"channel,omitempty"`
}

// LatestFilm contains the wanted information for the latest film
type LatestFilm struct {
	Title           string    `json:"title"`
	Link            string    `json:"link"`
	Rating          string    `json:"rating"`
	Year            string    `json:"year"`
	CreatedAt       time.Time `json:"created_at"`
	CreatedAtString string    `json:"created_at_string"`
}

// Collect returns the latest film in user's activity
// host: https://letterboxd.com
func (l *LatestFilm) Collect(host string, username string) error {
	resp, err := http.Get(fmt.Sprintf("%s/%s/rss", host, username))
	if err != nil {
		return errors.Wrap(err, "get films failed")
	}

	defer resp.Body.Close()

	var rss rssDocument
	err = xml.NewDecoder(resp.Body).Decode(&rss)
	if err != nil {
		return errors.Wrap(err, "body unmarshal failed")
	}

	if len(rss.Channel.Item) == 0 {
		return errors.New("there were no items in the feed")
	}

	itemTitle := rss.Channel.Item[0].Title.String
	regexTitle := regexp.MustCompile(`^(.*), `)
	regexYear := regexp.MustCompile(`^.*, (\d{4})`)
	regexRating := regexp.MustCompile(`^.*, \d{4} - (\S*)`)
	matchesTitle := regexTitle.FindStringSubmatch(itemTitle)
	matchesYear := regexYear.FindStringSubmatch(itemTitle)
	matchesRating := regexRating.FindStringSubmatch(itemTitle)

	if len(matchesTitle) != 2 {
		return fmt.Errorf("failed to get title from: '%v'", itemTitle)
	}
	if len(matchesYear) != 2 {
		return fmt.Errorf("failed to get year from: '%v'", itemTitle)
	}

	matchesString := ""
	if len(matchesRating) == 2 {
		matchesString = matchesRating[1]
	}

	createdAt, err := time.Parse("Mon, 2 Jan 2006 15:04:05 -0700", rss.Channel.Item[0].PubDate.String)
	if err != nil {
		return errors.Wrap(err, "failed to parse item date")
	}

	l.Title = matchesTitle[1]
	l.Year = matchesYear[1]
	l.Rating = matchesString
	l.CreatedAt = createdAt
	l.Link = rss.Channel.Item[0].Link.String

	return nil
}
