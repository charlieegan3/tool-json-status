package status

import (
	"time"

	"github.com/charlieegan3/tool-json-status/internal/pkg/collect"
	"github.com/charlieegan3/tool-json-status/internal/pkg/utils"
)

type Data struct {
	Tweet     collect.LatestTweet    `json:"tweet"`
	Post      collect.LatestPost     `json:"post"`
	Activity  collect.LatestActivity `json:"activity"`
	Film      collect.LatestFilm     `json:"film"`
	Commit    collect.LatestCommit   `json:"commit"`
	Play      collect.LatestPlay     `json:"play"`
	CreatedAt time.Time              `json:"created_at"`
}

func (d *Data) setCreatedAtStrings() {
	d.Tweet.CreatedAtString = utils.CompactHumanizeTime(d.Tweet.CreatedAt)
	d.Post.CreatedAtString = utils.CompactHumanizeTime(d.Post.CreatedAt)
	d.Activity.CreatedAtString = utils.CompactHumanizeTime(d.Activity.CreatedAt)
	d.Film.CreatedAtString = utils.CompactHumanizeTime(d.Film.CreatedAt)
	d.Commit.CreatedAtString = utils.CompactHumanizeTime(d.Commit.CreatedAt)
	d.Play.CreatedAtString = utils.CompactHumanizeTime(d.Play.CreatedAt)
}
