package binary

import (
	"regexp"

	"github.com/deifyed/xctl/pkg/tools/clients/kubectl"
)

func errorHandler(err error, defaultError error) error {
	switch {
	case isConnectionRefused(err):
		return kubectl.ErrConnectionRefused
	case isErrNotFound(err):
		return kubectl.ErrNotFound
	default:
		return defaultError
	}
}

var reErrNotFound = regexp.MustCompile(`Error from server \(NotFound\): pods ".+" not found\W`)

func isErrNotFound(err error) bool {
	return reErrNotFound.MatchString(err.Error())
}

var connectionRefusedRe = regexp.MustCompile(`.*connection refused.*`)

func isConnectionRefused(err error) bool {
	return connectionRefusedRe.MatchString(err.Error())
}
