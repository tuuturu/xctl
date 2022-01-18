package binary

import (
	"regexp"
)

var reConnectionRefused = regexp.MustCompile(`.*connect: connection refused.*`)

func isConnectionRefused(err error) bool {
	return reConnectionRefused.MatchString(err.Error())
}
