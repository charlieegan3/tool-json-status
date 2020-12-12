package collect

import (
	"encoding/base64"
	"encoding/json"
	"log"
	"time"

	"github.com/charlieegan3/json-charlieegan3/internal/pkg/proxy"
	"github.com/charlieegan3/json-charlieegan3/internal/pkg/types"
	"github.com/pkg/errors"
)

// LatestPost stores the URL, location and time of the latest post, video or photo
type LatestPost struct {
	URL             string    `json:"url"`
	Location        string    `json:"location"`
	Type            string    `json:"type"`
	CreatedAt       time.Time `json:"created_at"`
	CreatedAtString string    `json:"created_at_string"`
}

// Collect returns latest post for a given user
func (l *LatestPost) Collect(url, username, cookie string) error {
	posts, err := latestPosts(url, username, cookie)
	if err != nil {
		return errors.Wrap(err, "failed to get latest posts for user")
	}

	post := posts[0]
	postType := "photo"
	if post.IsVideo == true {
		postType = "video"
	}
	createdAt := time.Unix(post.TakenAtTimestamp, 0)

	l.Location = post.Location.Name
	l.Type = postType
	l.URL = "https://instagram.com/p/" + post.Shortcode
	l.CreatedAt = createdAt

	return nil
}

func latestPosts(url, username, cookie string) ([]types.LatestPost, error) {
	bytes, err := base64.StdEncoding.DecodeString(cookie)
	if err != nil {
		log.Fatal("INSTAGRAM_COOKIE_STRING must be b64")
	}
	cookie = string(bytes)
	var posts []types.LatestPost

	headers := map[string]string{
		"Cookie": cookie,
	}

	_, body, err := proxy.GetURLViaProxy(url+username+"/?__a=1", headers)
	if err != nil {
		return posts, errors.Wrap(err, "failed to get url via proxy")
	}

	var profile types.Profile
	if err := json.Unmarshal(body, &profile); err != nil {
		return posts, errors.Wrap(err, "failed to parse response")
	}

	for _, v := range profile.Graphql.User.EdgeOwnerToTimelineMedia.Edges {
		posts = append(posts, v.Node)
	}

	return posts, nil
}
