package cloud

import "github.com/pkg/errors"

var (
	// ErrNotFound indicates something is missing
	ErrNotFound = errors.New("not found")
	// ErrNotAuthenticated indicates invalid or missing authentication
	ErrNotAuthenticated = errors.New("not authenticated")
)
