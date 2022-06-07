package secrets

import "errors"

var (
	// ErrNotFound indicates something is missing
	ErrNotFound = errors.New("not found")
	// ErrUserAborted indicates user denied the operation
	ErrUserAborted = errors.New("user aborted")
)
