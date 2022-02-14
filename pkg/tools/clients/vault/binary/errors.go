package binary

import (
	"regexp"

	"github.com/deifyed/xctl/pkg/tools/clients/vault"
)

func errorHandler(err error, defaultError error) error {
	switch {
	case isConnectionRefused(err):
		return vault.ErrConnectionRefused
	default:
		return defaultError
	}
}

var reConnectionRefused = regexp.MustCompile(`.*connect: connection refused.*`)

func isConnectionRefused(err error) bool {
	return reConnectionRefused.MatchString(err.Error())
}
