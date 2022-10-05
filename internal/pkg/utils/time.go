package utils

import (
	"github.com/dustin/go-humanize"
	"strings"
	"time"
)

func CompactHumanizeTime(time time.Time) string {
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
