package cloud

import (
	"github.com/deifyed/xctl/pkg/tools/i18n"
	"github.com/pkg/errors"
)

var (
	// ErrNotFound indicates something is missing
	ErrNotFound = errors.New("not found")
	// ErrTimeout indicates something ran out of time
	ErrTimeout = errors.New("timeout")
	// ErrNotAuthenticated indicates invalid or missing authentication
	ErrNotAuthenticated = i18n.HumanReadableError{
		Content: "not authenticated",
		Key:     "cloud/notAuthenticated",
	}
)
