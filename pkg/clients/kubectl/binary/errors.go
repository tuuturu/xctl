package binary

import "regexp"

var reErrNotFound = regexp.MustCompile(`Error from server \(NotFound\): pods ".+" not found\s`)

func isErrNotFound(err error) bool {
	return reErrNotFound.MatchString(err.Error())
}

var connectionRefusedRe = regexp.MustCompile(`.*connection refused.*`)

func isConnectionRefused(s string) bool {
	return connectionRefusedRe.MatchString(s)
}
